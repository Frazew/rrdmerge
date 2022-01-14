package merger

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
						toCopy <- filePair{src: fileA, dst: fileB}
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
				merge(pair.src, pair.dst)
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

func merge(src string, dst string) {
	rrdA, err := rrdtool.Load(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load rrd file %s: %s\n", src, err)
		return
	}

	rrdB, err := rrdtool.Load(dst)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load rrd file %s: %s\n", dst, err)
		return
	}

	skipSteps := int(rrdB.Lastupdate-rrdA.Lastupdate) / 300
	if skipSteps > 0 {
		skipSteps--
	}

	if skipSteps < 0 {
		fmt.Fprintf(os.Stderr, "Failed to merge rrd files %s and %s: B has an older lastupdate value than A, maybe reverse them?\n", src, dst)
		return
	}

	if len(rrdA.Ds) != len(rrdB.Ds) {
		fmt.Fprintf(os.Stderr, "Failed to merge rrd files %s and %s: they don't have the same number of DS\n", src, dst)
		return
	}

	rraAMapping := make(map[rrdMapping]rrd.Rra, len(rrdA.Rra))

	for _, rraA := range rrdA.Rra {
		// We only support the CF_AVERAGE, CF_MAXIMUM, CF_MINIMUM, and CF_LAST functions when unmarshalling
		rraAMapping[rrdMapping{rraA.Cf, rraA.PdpPerRow, rraA.Params.Xff}] = rraA
	}
	for _, rraB := range rrdB.Rra {
		if rraA, ok := rraAMapping[rrdMapping{rraB.Cf, rraB.PdpPerRow, rraB.Params.Xff}]; ok {
			rowALength := len(rraA.Database.Row)
			for i, row := range rraB.Database.Row {
				if i+skipSteps/rraA.PdpPerRow >= rowALength {
					break
				}
				for j, v := range row.V {
					if v == "NaN" && rraA.Database.Row[i].V[j] != "NaN" { // If we have a value, let's not override it with NaN
						row.V[j] = rraA.Database.Row[i+skipSteps/rraA.PdpPerRow].V[j]
					}
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Failed to find a match for RRA %s,%d,%s in %s\n", rraB.Cf, rraB.PdpPerRow, rraB.Params.Xff, src)
		}
	}

	info, err := os.Stat(dst)
	err = rrdtool.Dump(rrdB, dst, info.Mode())
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
