package merger

import (
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

func TestMerge_LostHistory(t *testing.T) {
	copyFile("../../fixtures/test_lost_history/b.rrd", "../../fixtures/test_lost_history/b_tmp.rrd")
	merge("../../fixtures/test_lost_history/a.rrd", "../../fixtures/test_lost_history/b_tmp.rrd", "", "", false)

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
	merge("../../fixtures/test_load.rrd", "../../fixtures/test_load_tmp.rrd", "", "", false)

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
