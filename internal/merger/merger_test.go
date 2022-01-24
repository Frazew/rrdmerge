package merger

import (
	"io/fs"
	"path"
	"path/filepath"
	"testing"

	"github.com/scaleway/rrdmerge/internal/rrd"
)

func TestLoadRrd(t *testing.T) {
	rrdPtr, err := loadRrd("../../fixtures/test_load.rrd")
	if err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}
	if rrdPtr == nil {
		t.Errorf("Expected rrdPtr not to be nil")
		return
	}
	if rrdPtr.Header.Version != "0003" {
		t.Errorf("Expected Version to be 0003, got %s", rrdPtr.Header.Version)
	}
	if len(rrdPtr.RraStore) != int(rrdPtr.Header.RraCount) {
		t.Errorf("Expected the RRA store to have %d entries, got %d", rrdPtr.Header.RraCount, len(rrdPtr.RraStore))
	}
	if len(rrdPtr.RraDataStore) != int(rrdPtr.Header.RraCount) {
		t.Errorf("Expected the RRA data store to have %d entries, got %d", rrdPtr.Header.RraCount, len(rrdPtr.RraDataStore))
	}
}

func TestIsValidFile(t *testing.T) {
	mergeSpec := MergeSpec{
		NoSkip: false,
	}
	_, err := mergeSpec.isValidFile("../../fixtures/test_load.rrd")
	if err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}

	_, err = mergeSpec.isValidFile("../../fixtures/")
	if err == nil {
		t.Errorf("Expected err not to be nil")
		return
	}
}

func TestCopyFile(t *testing.T) {
	mergeSpec := MergeSpec{
		NoSkip: false,
	}
	copyFile("../../fixtures/test_load.rrd", "../../fixtures/test_load_tmp.rrd")
	if _, err := mergeSpec.isValidFile("../../fixtures/test_load_tmp.rrd"); err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}
}

func TestMerge_LostHistory(t *testing.T) {
	copyFile("../../fixtures/test_lost_history/b.rrd", "../../fixtures/test_lost_history/b_tmp.rrd")
	mergeSpec := MergeSpec{
		RrdA:        "../../fixtures/test_lost_history/a.rrd",
		RrdB:        "../../fixtures/test_lost_history/b_tmp.rrd",
		MergeType:   MergeFile,
		Concurrency: 4,
		Common:      false,
		NoSkip:      false,
		DaemonOpt:   "",
		StripPath:   "",
		DryRun:      false,
	}
	mergeSpec.DoMerge()

	testLostHistoryTarget, err := loadRrd("../../fixtures/test_lost_history/target.rrd")
	if err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}
	testLostHistoryB, err := loadRrd("../../fixtures/test_lost_history/b_tmp.rrd")
	if err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}

	if err = rrd.FromRRDStruct(testLostHistoryTarget).Equals(*rrd.FromRRDStruct(testLostHistoryB)); err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}
}

func TestMerge_SameFile(t *testing.T) {
	copyFile("../../fixtures/test_load.rrd", "../../fixtures/test_load_tmp.rrd")
	mergeSpec := MergeSpec{
		RrdA:        "../../fixtures/test_load.rrd",
		RrdB:        "../../fixtures/test_load_tmp.rrd",
		MergeType:   MergeFile,
		Concurrency: 4,
		Common:      false,
		NoSkip:      false,
		DaemonOpt:   "",
		StripPath:   "",
		DryRun:      false,
	}
	mergeSpec.DoMerge()

	testLoad, err := loadRrd("../../fixtures/test_load.rrd")
	if err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}
	testLoadMerged, err := loadRrd("../../fixtures/test_load_tmp.rrd")
	if err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}

	if err = rrd.FromRRDStruct(testLoad).Equals(*rrd.FromRRDStruct(testLoadMerged)); err != nil {
		t.Errorf("Expected err to be nil, got: %w", err)
		return
	}
}

func TestMerge_SameFolder(t *testing.T) {
	mergeSpec := MergeSpec{
		RrdA:        "../../fixtures/merge_into/",
		RrdB:        "../../fixtures/merge_into_tmp/",
		MergeType:   MergeFolder,
		Concurrency: 4,
		Common:      false,
		NoSkip:      false,
		DaemonOpt:   "",
		StripPath:   "",
		DryRun:      false,
	}
	mergeSpec.DoMerge()

	err := filepath.WalkDir("../../fixtures/merge_into", func(file string, d fs.DirEntry, err error) error {
		_, err = mergeSpec.isValidFile(file)
		if err == nil {
			testLoad, err := loadRrd(file)
			if err != nil {
				t.Errorf("Expected err to be nil, got: %w", err)
				return err
			}
			testLoadMerged, err := loadRrd(path.Join(mergeSpec.RrdB, path.Base(file)))
			if err != nil {
				t.Errorf("Expected err to be nil, got: %w", err)
				return err
			}

			if err = rrd.FromRRDStruct(testLoad).Equals(*rrd.FromRRDStruct(testLoadMerged)); err != nil {
				t.Errorf("Expected err to be nil, got: %w", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		t.Errorf(err.Error())
		return
	}

}
