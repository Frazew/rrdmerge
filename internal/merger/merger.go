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
}

func (spec MergeSpec) String() string {
	return fmt.Sprintf("a: %s, b: %s, common: %t, type: %s", spec.RrdA, spec.RrdB, spec.Common, spec.MergeType)
}

func (spec MergeSpec) DoMerge() error {
	fmt.Println(spec)
	if spec.MergeType == MergeFolder {
		return spec.mergeFolder()
	} else {
		if spec.Common {
			fmt.Println("Ignoring flag common because we're not merging folders")
		}
		return spec.mergeFile()
	}
}

func (spec MergeSpec) mergeFolder() error {
	var wgWalker sync.WaitGroup
	var wgCopier sync.WaitGroup
	var wgMerger sync.WaitGroup
	files := make(chan string)
	toCopy := make(chan filePair)
	toMerge := make(chan filePair)

	go func() {
		defer close(files)

		filepath.WalkDir(spec.RrdA, func(file string, d fs.DirEntry, err error) error {
			if err == nil {
				if file, err := spec.isValidFile(file); err == nil {
					files <- file
				} else {
					fmt.Fprintf(os.Stderr, "Skipping file %s: %s\n", file, err)
				}
			} else {
				fmt.Fprintf(os.Stderr, "Failed to read %s: %s\n", file, err)
			}
			return nil
		})
	}()

	for i := 0; i < spec.Concurrency; i++ {
		wgWalker.Add(1)
		go func() error {
			defer wgWalker.Done()
			for fileA := range files {
				info, err := os.Stat(fileA)
				if err != nil {
					panic(err) // This should never happen since we've already stat'd the file before
				} else {
					fileB := filepath.Join(spec.RrdB, info.Name())
					if _, err := os.Stat(fileB); err == nil {
						if fileB, err := spec.isValidFile(fileB); err == nil {
							fmt.Fprintf(os.Stderr, "Merging files %s and %s\n", fileA, fileB)
							toMerge <- filePair{src: fileA, dst: fileB}
						} else {
							fmt.Fprintf(os.Stderr, "Not merging %s into %s because the target is not valid\n", fileA, fileB)
							continue
						}
					} else { // Simple copy since the target does not exist
						if !spec.Common {
							toCopy <- filePair{src: fileA, dst: fileB}
						}
					}
				}
			}
			return nil
		}()

		wgCopier.Add(1)
		go func() {
			defer wgCopier.Done()
			for pair := range toCopy {
				copyFile(pair.src, pair.dst)
			}
		}()
		wgMerger.Add(1)
		go func() {
			defer wgMerger.Done()
			for pair := range toMerge {
				merge(pair.src, pair.dst, spec)
			}
		}()
	}

	wgWalker.Wait()
	close(toCopy)
	close(toMerge)
	wgCopier.Wait()
	wgMerger.Wait()
	return nil
}

func (spec MergeSpec) mergeFile() error {
	merge(spec.RrdA, spec.RrdB, spec)
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
	} else {
		if info.IsDir() {
			return file, errors.New("path is a directory")
		} else {
			return file, errors.New("file is non-regular or doesn't have a .rrd extension (try the noSkip flag?)")
		}
	}
}

func loadRrd(path string) (*rrd.Rrd, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	rrdPtr := rrd.NewRrd()
	return rrdPtr, rrdPtr.Read(kaitai.NewStream(file), rrdPtr, rrdPtr)
}

func merge(src string, dst string, spec MergeSpec) {
	if spec.DaemonOpt != "" {
		err := rrdtool.Tune(dst, spec.DaemonOpt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to flush rrd file %s: %s\n", src, err)
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

	// We merge into the newest rrd
	if rrdA.LiveHead.LastUpdate > rrdB.LiveHead.LastUpdate {
		rrdA, rrdB = rrdB, rrdA
	}

	stepsDifference := int(rrdB.LiveHead.LastUpdate-rrdA.LiveHead.LastUpdate) / int(rrdA.Header.PdpStep)

	for rraIdx, rra := range rrdB.RraDataStore {
		if stepsDifference/int(rrdB.RraStore[rraIdx].PdpCount) > int(rra.RowCount) {
			fmt.Fprintf(os.Stderr, "Not merging data in RRA %d in %s because it has already rolled over\n", rraIdx, src)
			continue
		}

		timeShift := int(rrdB.LiveHead.LastUpdate-rrdA.LiveHead.LastUpdate)/int(rrdA.Header.PdpStep*rrdA.RraStore[rraIdx].PdpCount) + 1

		k := 0
		startFrom := int(rrdA.RraPtrStore[rraIdx])
		oldIndex := startFrom - k

		totalRecovered := 0
		for i := int(rrdB.RraPtrStore[rraIdx]) - timeShift; i >= 0; i-- {
			oldIndex = startFrom - k

			if totalRecovered >= int(rrdA.RraStore[rraIdx].RowCount) {
				break
			}

			for j := range rra.Row[i].Values {
				if !math.IsNaN(rrdA.RraDataStore[rraIdx].Row[oldIndex].Values[j]) {
					rrdB.RraDataStore[rraIdx].Row[i].Values[j] = rrdA.RraDataStore[rraIdx].Row[oldIndex].Values[j]
				}
			}
			k++
			totalRecovered++
			if oldIndex == 0 {
				k = 1
				startFrom = int(rrdA.RraStore[rraIdx].RowCount)
			}
		}

		for i := int(rra.RowCount) - 1; i > int(rrdB.RraPtrStore[rraIdx]); i-- {
			oldIndex := startFrom - k

			if totalRecovered >= int(rrdA.RraStore[rraIdx].RowCount) {
				break
			}

			for j := range rra.Row[i].Values {
				if !math.IsNaN(rrdA.RraDataStore[rraIdx].Row[oldIndex].Values[j]) {
					rrdB.RraDataStore[rraIdx].Row[i].Values[j] = rrdA.RraDataStore[rraIdx].Row[oldIndex].Values[j]
				}
			}
			k++
			totalRecovered++
			if oldIndex == 0 {
				k = 1
				startFrom = int(rrdA.RraStore[rraIdx].RowCount)
			}
		}
	}

	info, err := os.Stat(dst)
	err = rrdtool.Restore(rrdB, dst, info.Mode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write destination rrd file %s: %s\n", dst, err)
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
