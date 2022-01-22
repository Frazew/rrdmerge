package main

import (
	"testing"

	"github.com/scaleway/rrdmerge/internal/merger"
)

func TestPathToMergeType(t *testing.T) {
	pathToSpec := map[string]merger.MergeType{
		"../../fixtures":               merger.MergeFolder,
		"../../fixtures/test_load.rrd": merger.MergeFile,
	}

	for path, spec := range pathToSpec {
		mergeType, fileInfo, err := pathToMergeType(path)
		if err != nil {
			t.Errorf("Expected err to be nil, got: %w", err)
			return
		}
		if fileInfo == nil {
			t.Errorf("Expected fileInfo not to be nil")
			return
		}
		if mergeType != spec {
			t.Errorf("Expected mergeType to be %d, got %d", spec, mergeType)
			return
		}
	}
}
