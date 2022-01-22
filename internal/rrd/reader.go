package rrd

import (
	"errors"
	"fmt"
)

type RRDReader struct {
	rrdPtr   *Rrd
	curRra   *Rrd_RraData
	startIdx int
	endIdx   int
	headIdx  int
	read     int
}

func Reader(rrdPtr *Rrd) *RRDReader {
	reader := new(RRDReader)
	reader.rrdPtr = rrdPtr
	reader.SelectRRA(0)
	return reader
}

func (reader *RRDReader) SelectRRA(index int) error {
	if index > int(reader.rrdPtr.Header.RraCount) {
		return errors.New(fmt.Sprintf("RRA at index %d does not exist", index))
	}
	reader.curRra = reader.rrdPtr.RraDataStore[index]
	reader.startIdx = int(reader.rrdPtr.RraPtrStore[index]) + 1
	reader.endIdx = int(reader.rrdPtr.RraPtrStore[index])
	reader.read = 0
	reader.Seek(0)
	return nil
}

func (reader *RRDReader) Next() (*Rrd_RraDataRow, int, int, error) {
	if reader.headIdx == reader.endIdx {
		return reader.curRra.Row[reader.headIdx], reader.read, reader.headIdx, errors.New("Reached end of RRA")
	}

	defer func() {
		reader.headIdx++
		reader.read++

		if reader.headIdx == int(reader.curRra.RowCount) {
			reader.headIdx = 0
		}
	}()
	return reader.curRra.Row[reader.headIdx], reader.read, reader.headIdx, nil
}

func (reader *RRDReader) Seek(position int) error {
	if position >= int(reader.curRra.RowCount) {
		return errors.New("Seeking at an index higer than the row count")
	}
	reader.headIdx = (reader.startIdx + position) % int(reader.curRra.RowCount)
	return nil
}
