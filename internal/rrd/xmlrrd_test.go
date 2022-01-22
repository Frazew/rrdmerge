package rrd

import (
	"encoding/xml"
	"os"
	"testing"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func TestFromRRDStruct(t *testing.T) {
	rrdA, err := os.Open("../../fixtures/test_load.rrd")
	if err != nil {
		t.Errorf("Failed to open fixture test_load.rrd")
	}

	xmlABytes, err := os.ReadFile("../../fixtures/test_load.xml")
	if err != nil {
		t.Errorf("Failed to open fixture test_load.xml")
	}
	var xmlRrd xmlRrd
	err = xml.Unmarshal(xmlABytes, &xmlRrd)
	if err != nil {
		t.Errorf("Expected err not to be nil, got: %w", err)
	}

	rrdPtr := NewRrd()
	rrdPtr.Read(kaitai.NewStream(rrdA), rrdPtr, rrdPtr)
	if rrdPtr == nil {
		t.Errorf("Expected rrdPtr not to be nil")
	}
	xmlOutput := FromRRDStruct(rrdPtr)

	err = xmlOutput.Equals(xmlRrd)
	if err != nil {
		t.Errorf("Expected err not to be nil, got: %w", err)
	}
}
