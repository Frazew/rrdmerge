package merger

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"

	"github.com/scaleway/rrdmerge/internal/rrd"
	"github.com/scaleway/rrdmerge/internal/rrdtool"
)

type filePair struct {
	src, dst string
}

type rrdMapping struct {
	DF  string
	PdP int
	Xff string
}

type MergeType int

const (
	MergeFolder MergeType = iota
	MergeFile
)

func (mergeType MergeType) String() string {
	switch mergeType {
	case MergeFolder:
		return "folder"
	case MergeFile:
		return "file"
	default:
		return "unknown"
	}
}

type MergeSpec struct {
	RrdA        string
	RrdB        string
	MergeType   MergeType
	Concurrency int
	Common      bool
	NoSkip      bool
	DaemonOpt   string
	StripPath   string
	DryRun      bool
}

func (spec MergeSpec) String() string {
	return fmt.Sprintf("a: %s, b: %s, common: %t, type: %s", spec.RrdA, spec.RrdB, spec.Common, spec.MergeType)
}

func (spec MergeSpec) DoMerge() error {
	if spec.MergeType == MergeFolder {
		return spec.mergeFolder()
	}

	if spec.Common {
		fmt.Println("Ignoring flag common because we're not merging folders")
	}
	return spec.mergeFile()
}

func (spec MergeSpec) mergeFolder() error {
	var wgCopier sync.WaitGroup
	var wgMerger sync.WaitGroup
	toCopy := make(chan *filePair, 100)
	toMerge := make(chan *filePair, 100)

	for i := 0; i < spec.Concurrency; i++ {
		wgCopier.Add(1)
		go func(c chan *filePair, dryRun bool) {
			defer wgCopier.Done()
			for pair := range c {
				fmt.Fprintf(os.Stderr, "Copying %s to %s\n", pair.src, pair.dst)
				if !dryRun {
					copyFile(pair.src, pair.dst)
				}
			}
		}(toCopy, spec.DryRun)

		wgMerger.Add(1)
		go func(c chan *filePair, daemonOpt string, stripPath string, dryRun bool) {
			defer wgMerger.Done()
			for pair := range c {
				fmt.Fprintf(os.Stderr, "Merging files %s and %s\n", pair.src, pair.dst)
				merge(pair.src, pair.dst, daemonOpt, stripPath, dryRun)
			}
		}(toMerge, spec.DaemonOpt, spec.StripPath, spec.DryRun)
	}

	filepath.WalkDir(spec.RrdA, func(file string, d fs.DirEntry, err error) error {
		if err == nil {
			if fileA, err := spec.isValidFile(file); err == nil {
				info, err := os.Stat(fileA)
				if err != nil {
					panic(err) // This should never happen since we've already stat'd the file before
				} else {
					fileB := filepath.Join(spec.RrdB, info.Name())
					if _, err := os.Stat(fileB); err == nil {
						if fileB, err := spec.isValidFile(fileB); err == nil {
							toMerge <- &filePair{src: fileA, dst: fileB}
						} else {
							fmt.Fprintf(os.Stderr, "Not merging %s into %s because the target is not valid\n", fileA, fileB)
							return nil
						}
					} else { // Simple copy since the target does not exist
						if !spec.Common {
							toCopy <- &filePair{src: fileA, dst: fileB}
						}
					}
				}
			} else {
				fmt.Fprintf(os.Stderr, "Skipping file %s: %s\n", file, err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %s\n", file, err)
		}
		return nil
	})
	close(toCopy)
	close(toMerge)

	wgCopier.Wait()
	wgMerger.Wait()
	return nil
}

func (spec MergeSpec) mergeFile() error {
	merge(spec.RrdA, spec.RrdB, spec.DaemonOpt, spec.StripPath, spec.DryRun)
	return nil
}

func (spec MergeSpec) isValidFile(file string) (string, error) {
	// Try to resolve symlinks
	info, err := os.Stat(file)
	if err != nil {
		return file, err
	}
	if info.Mode().Type() == os.ModeSymlink {
		file, err := filepath.EvalSymlinks(file)
		if err != nil {
			return file, err
		}
		info, err = os.Stat(file)
		if err != nil {
			return file, err
		}
	}
	if info.Mode().IsRegular() && (strings.HasSuffix(info.Name(), ".rrd") || spec.NoSkip) {
		return file, nil
	}
	if info.IsDir() {
		return file, errors.New("path is a directory")
	}
	return file, errors.New("file is non-regular or doesn't have a .rrd extension (try the noSkip flag?)")
}

func loadRrd(path string) (*rrd.Rrd, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	rrdPtr := rrd.NewRrd()
	return rrdPtr, rrdPtr.Read(kaitai.NewStream(file), rrdPtr, rrdPtr)
}

func merge(src string, dst string, daemonOpt string, stripPath string, dryRun bool) {
	if daemonOpt != "" {
		err := rrdtool.Flush(strings.Replace(dst, stripPath, "", 1), daemonOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to flush rrd file %s: %s\n", dst, err)
			return
		}
	}
	rrdA, err := loadRrd(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load rrd file %s: %s\n", src, err)
		return
	}

	rrdB, err := loadRrd(dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load rrd file %s: %s\n", dst, err)
		return
	}

	if rrdA.Header.DsCount != rrdB.Header.DsCount {
		fmt.Fprintf(os.Stderr, "Failed to merge rrd files %s and %s: they don't have the same number of DS\n", src, dst)
		return
	}

	if rrdB.LiveHead.LastUpdate-rrdA.LiveHead.LastUpdate == 0 {
		fmt.Fprintf(os.Stderr, "Not merging rrd files %s and %s: they have the same last update value\n", src, dst)
		return
	}

	// We merge into the newest rrd
	if rrdA.LiveHead.LastUpdate > rrdB.LiveHead.LastUpdate {
		rrdA, rrdB = rrdB, rrdA
	}

	lastUpdateDifference := int(rrdB.LiveHead.LastUpdate - rrdA.LiveHead.LastUpdate)
	stepsDifference := (lastUpdateDifference + int(rrdA.Header.PdpStep) - 1) / int(rrdA.Header.PdpStep)

	rrdAReader := rrd.Reader(rrdA)
	rrdBReader := rrd.Reader(rrdB)
	for rraIdx, rra := range rrdB.RraDataStore {
		if stepsDifference/int(rrdB.RraStore[rraIdx].PdpCount) > int(rra.RowCount) {
			continue
		}
		if rraIdx > int(rrdA.Header.RraCount)-1 {
			break
		}

		stepSize := int(rrdB.RraStore[rraIdx].PdpCount) * int(rrdA.Header.PdpStep)
		rrdAStep := int(rrdA.LiveHead.LastUpdate) / stepSize
		rrdBStep := int(rrdB.LiveHead.LastUpdate) / stepSize
		timeShift := rrdBStep - rrdAStep

		rrdAReader.SelectRRA(rraIdx)
		rrdBReader.SelectRRA(rraIdx)
		rrdAReader.Seek(timeShift)

		for {
			rrdBRow, _, _, err := rrdBReader.Next()
			if err != nil {
				break
			}

			rrdARow, _, _, err := rrdAReader.Next()

			for j := range rrdBRow.Values {
				if !math.IsNaN(rrdARow.Values[j]) {
					rrdBRow.Values[j] = rrdARow.Values[j]
				}
			}
			if err != nil {
				break
			}
		}
	}

	if !dryRun {
		err = rrdtool.Restore(rrdB, dst)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write destination rrd file %s: %s\n", dst, err)
		}
	}
}

func copyFile(src string, dst string) {
	fileAReader, err := os.Open(src)
	defer fileAReader.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s for reading: %s\n", src, err)
		return
	}

	fileBWriter, err := os.Create(dst)
	defer fileBWriter.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %s for writing: %s\n", dst, err)
		return
	}

	_, err = io.Copy(fileBWriter, fileAReader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing from %s to %s: %s\n", src, dst, err)
		return
	}

	err = fileBWriter.Sync()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while syncing file %s: %s\n", dst, err)
		return
	}
}
