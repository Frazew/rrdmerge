package rrd

import (
	"errors"
	"fmt"
)

// Reader is a helper to navigate within the RRD structure
type Reader struct {
	rrdPtr   *Rrd
	curRra   *Rrd_RraData
	startIdx int
	endIdx   int
	headIdx  int
	read     int
}

// NewReader creates a new Reader
func NewReader(rrdPtr *Rrd) *Reader {
	reader := new(Reader)
	reader.rrdPtr = rrdPtr
	reader.SelectRRA(0)
	return reader
}

// SelectRRA selects the RRA specified by the given index and initializes the Reader
func (reader *Reader) SelectRRA(index int) error {
	if index > int(reader.rrdPtr.Header.RraCount) {
		return fmt.Errorf("RRA at index %d does not exist", index)
	}
	reader.curRra = reader.rrdPtr.RraDataStore[index]
	reader.startIdx = int(reader.rrdPtr.RraPtrStore[index]) + 1
	reader.endIdx = int(reader.rrdPtr.RraPtrStore[index])
	reader.read = 0
	reader.Seek(0)
	return nil
}

// Next returns the next entry in the currently selected RRA and updates the internal state accordingly
func (reader *Reader) Next() (*Rrd_RraDataRow, error) {
	if reader.headIdx == reader.endIdx {
		return reader.curRra.Row[reader.headIdx], errors.New("Reached end of RRA")
	}

	defer func() {
		reader.headIdx++
		reader.read++

		if reader.headIdx == int(reader.curRra.RowCount) {
			reader.headIdx = 0
		}
	}()
	return reader.curRra.Row[reader.headIdx], nil
}

// Seek moves the current reader pointer to the specified position in the RRA
func (reader *Reader) Seek(position int) error {
	if position >= int(reader.curRra.RowCount) {
		return errors.New("Seeking at an index higer than the row count")
	}
	reader.headIdx = (reader.startIdx + position) % int(reader.curRra.RowCount)
	return nil
}
