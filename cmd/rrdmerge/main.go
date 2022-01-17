package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/scaleway/rrdmerge/internal/merger"
)

func pathToMergeType(path string) (merger.MergeType, fs.FileInfo, error) {
	if file, err := os.Stat(path); err == nil {
		if file.IsDir() {
			return merger.MergeFolder, file, nil
		}
		return merger.MergeFile, file, nil
	} else if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
		return -1, nil, err
	} else {
		return -1, nil, fmt.Errorf("Path %s could not be recognized", path)
	}
}

func main() {
	rrdA := flag.String("a", "", "The first file/folder to merge from (required)")
	rrdB := flag.String("b", "", "The second file/folder to merge to (required)")
	common := flag.Bool("common", false, "Merge only files that are common to both sources when merging folders (optional)")
	noSkip := flag.Bool("noSkip", false, "Do not skip files with an extension other that .rrd (optional)")
	daemonOpt := flag.String("d", "", "Flush the rrd files with the given rrdcached daemon before merging (optional)")
	concurrency := flag.Int("t", 4, "Run this many parallel merger jobs when processing a directory (optional)")

	flag.Parse()
	if flag.Parsed() {
		if *rrdA == "" || *rrdB == "" {
			flag.PrintDefaults()
			os.Exit(1)
		} else {
			if mergeTypeA, fileInfoA, err := pathToMergeType(*rrdA); err == nil {
				if mergeTypeB, fileInfoB, err := pathToMergeType(*rrdB); err == nil {
					if mergeTypeA != mergeTypeB {
						fmt.Fprintf(os.Stderr, "Incompatible path specifications between a: %s and b: %s\n", mergeTypeA, mergeTypeB)
						os.Exit(1)
					} else if os.SameFile(fileInfoA, fileInfoB) {
						fmt.Fprintf(os.Stderr, "Files a: %s and b: %s are the same\n", *rrdA, *rrdB)
						os.Exit(1)
					} else {
						mergeSpec := &merger.MergeSpec{
							RrdA:        *rrdA,
							RrdB:        *rrdB,
							MergeType:   mergeTypeA,
							Concurrency: *concurrency,
							Common:      *common,
							NoSkip:      *noSkip,
							DaemonOpt:   *daemonOpt,
						}
						start := time.Now()
						mergeSpec.DoMerge()
						elapsed := time.Since(start)
						fmt.Fprintf(os.Stderr, "Merged %s and %s in %s\n", *rrdA, *rrdB, elapsed)
					}
				} else {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
