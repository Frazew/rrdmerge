package rrd

import (
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
)

type xmlDsParams struct {
	PrimaryValue      string `xml:"primary_value"`
	SecondaryValue    string `xml:"secondary_value"`
	Value             string `xml:"value"`
	UnknownDatapoints string `xml:"unknown_datapoints"`
}

type xmlRraParams struct {
	Xff string `xml:"xff"` // Only support CF_AVERAGE, CF_MAXIMUM, CF_MINIMUM, and CF_LAST
}

type xmlRraCdpPrep struct {
	Ds []xmlDsParams `xml:"ds"`
}

type xmlRraRow struct {
	V []string `xml:"v"`
}

type xmlRraDatabase struct {
	Rows []xmlRraRow `xml:"row"`
}

type xmlRra struct {
	Cf        string         `xml:"cf"`
	PdpPerRow int            `xml:"pdp_per_row"`
	Params    xmlRraParams   `xml:"params"`
	CdpPrep   xmlRraCdpPrep  `xml:"cdp_prep"`
	Database  xmlRraDatabase `xml:"database"`
}

type xmlDs struct {
	Name             string `xml:"name"`
	Type             string `xml:"type"`
	MinimalHeartbeat string `xml:"minimal_heartbeat"`
	Min              string `xml:"min"`
	Max              string `xml:"max"`
	LastDs           string `xml:"last_ds"`
	Value            string `xml:"value"`
	UnknownSec       string `xml:"unknown_sec"`
}

type xmlRrd struct {
	XMLName    xml.Name `xml:"rrd"`
	Version    string   `xml:"version"`
	Step       int      `xml:"step"`
	Lastupdate int64    `xml:"lastupdate"`
	Ds         []xmlDs  `xml:"ds"`
	Rra        []xmlRra `xml:"rra"`
}

func FromRRDStruct(rrdPtr *Rrd) *xmlRrd {
	dsSlice := make([]xmlDs, rrdPtr.Header.DsCount)
	for dsIdx, ds := range rrdPtr.DsStore {
		dsSlice[dsIdx] = xmlDs{
			Name:             ds.Name,
			Type:             ds.Datasource,
			MinimalHeartbeat: fmt.Sprintf("%d", ds.Params.MinHeartbeatCount),
			Min:              fmt.Sprintf("%.10e", ds.Params.MinVal),
			Max:              fmt.Sprintf("%.10e", ds.Params.MaxVal),
			LastDs:           rrdPtr.PdpPrepStore[dsIdx].LastDsReading,
			Value:            fmt.Sprintf("%.10e", rrdPtr.PdpPrepStore[dsIdx].Params.CurrentValue),
			UnknownSec:       fmt.Sprintf("%d", rrdPtr.PdpPrepStore[dsIdx].Params.UnknownSecCount),
		}
	}

	rraSlice := make([]xmlRra, rrdPtr.Header.RraCount)
	for rraIdx, rra := range rrdPtr.RraDataStore {
		dsParams := make([]xmlDsParams, rrdPtr.Header.DsCount)
		for cdpIdx, cdp := range rrdPtr.CdpPrepStore[rraIdx].(*Rrd_CdpPrep).Params {
			dsParams[cdpIdx] = xmlDsParams{
				PrimaryValue:      fmt.Sprintf("%.10e", cdp.PrimaryValue),
				SecondaryValue:    fmt.Sprintf("%.10e", cdp.SecondaryValue),
				Value:             fmt.Sprintf("%.10e", cdp.Value),
				UnknownDatapoints: fmt.Sprintf("%d", cdp.UnknownPdpCount),
			}
		}

		rows := make([]xmlRraRow, rra.RowCount)

		written := 0
		position := int(rrdPtr.RraPtrStore[rraIdx]) + 1
		for written < int(rra.RowCount) {
			rows[written].V = make([]string, rrdPtr.Header.DsCount)

			if position >= int(rra.RowCount) {
				position = 0
			}
			for i, v := range rra.Row[position].Values {
				rows[written].V[i] = fmt.Sprintf("%.10e", v)
			}

			position++
			written++
		}

		rraSlice[rraIdx] = xmlRra{
			Cf:        rrdPtr.RraStore[rraIdx].Cf,
			PdpPerRow: int(rrdPtr.RraStore[rraIdx].PdpCount),
			Params:    xmlRraParams{fmt.Sprintf("%.10e", rrdPtr.RraStore[rraIdx].Params.(*Rrd_RraParams).Xff)},
			CdpPrep:   xmlRraCdpPrep{Ds: dsParams},
			Database:  xmlRraDatabase{Rows: rows},
		}
	}

	return &xmlRrd{
		Version:    "0003",
		Step:       int(rrdPtr.Header.PdpStep),
		Lastupdate: int64(rrdPtr.LiveHead.LastUpdate),
		Ds:         dsSlice,
		Rra:        rraSlice,
	}
}

func (xmlA xmlRrd) Equals(xmlB xmlRrd) error {
	if !reflect.DeepEqual(xmlA.Lastupdate, xmlB.Lastupdate) {
		return errors.New("Expected the last update value to be equal")
	}

	if !reflect.DeepEqual(xmlA.Step, xmlB.Step) {
		return errors.New("Expected the step size to be equal")
	}

	if len(xmlA.Rra) != len(xmlB.Rra) {
		return errors.New("Expected to have the same count of RRAs")
	}

	if len(xmlA.Ds) != len(xmlB.Ds) {
		return errors.New("Expected to have the same count of DSs")
	}

	for dsIdx, ds := range xmlB.Ds {
		if ds.Name != xmlA.Ds[dsIdx].Name {
			return errors.New("Expected to have the same name")
		}
		if ds.Type != xmlA.Ds[dsIdx].Type {
			return errors.New("Expected to have the same type")
		}
		if ds.MinimalHeartbeat != xmlA.Ds[dsIdx].MinimalHeartbeat {
			return errors.New("Expected to have the same heartbeat")
		}
		if !floatEquals(ds.Min, xmlA.Ds[dsIdx].Min) {
			return errors.New("Expected to have the same minimum")
		}
		if !floatEquals(ds.Max, xmlA.Ds[dsIdx].Max) {
			return errors.New("Expected to have the same maximum")
		}
		if ds.LastDs != xmlA.Ds[dsIdx].LastDs {
			return errors.New("Expected to have the same last DS value")
		}
		if !floatEquals(ds.Value, xmlA.Ds[dsIdx].Value) {
			return errors.New("Expected to have the same value")
		}
		if ds.UnknownSec != xmlA.Ds[dsIdx].UnknownSec {
			return errors.New("Expected to have the same count of unknown seconds")
		}
	}

	for rraIdx, rra := range xmlB.Rra {
		if !reflect.DeepEqual(rra.Params, xmlA.Rra[rraIdx].Params) {
			return errors.New(fmt.Sprintf("Expected RRA %d to have the same params field", rraIdx))
		}
		if rra.Cf != xmlA.Rra[rraIdx].Cf {
			return errors.New(fmt.Sprintf("Expected RRA %d to have the same CF", rraIdx))
		}
		if rra.PdpPerRow != xmlA.Rra[rraIdx].PdpPerRow {
			return errors.New(fmt.Sprintf("Expected RRA %d to have the same number of pdp per row", rraIdx))
		}

		if len(rra.CdpPrep.Ds) != len(xmlA.Rra[rraIdx].CdpPrep.Ds) {
			return errors.New(fmt.Sprintf("Expected RRA %d to have the same number of DS", rraIdx))
		}

		for dsIdx, ds := range rra.CdpPrep.Ds {
			if !floatEquals(ds.PrimaryValue, xmlA.Rra[rraIdx].CdpPrep.Ds[dsIdx].PrimaryValue) {
				return errors.New(fmt.Sprintf("Expected RRA %d DS %d to have the same primary value", rraIdx, dsIdx))
			}
			if !floatEquals(ds.SecondaryValue, xmlA.Rra[rraIdx].CdpPrep.Ds[dsIdx].SecondaryValue) {
				return errors.New(fmt.Sprintf("Expected RRA %d DS %d to have the same secondary value", rraIdx, dsIdx))
			}
			if !floatEquals(ds.Value, xmlA.Rra[rraIdx].CdpPrep.Ds[dsIdx].Value) {
				return errors.New(fmt.Sprintf("Expected RRA %d DS %d to have the same value", rraIdx, dsIdx))
			}
			if ds.UnknownDatapoints != xmlA.Rra[rraIdx].CdpPrep.Ds[dsIdx].UnknownDatapoints {
				return errors.New(fmt.Sprintf("Expected RRA %d DS %d to have the count of unknown datapoints", rraIdx, dsIdx))
			}
		}

		if len(rra.Database.Rows) != len(xmlA.Rra[rraIdx].Database.Rows) {
			return errors.New(fmt.Sprintf("Expected RRA %d to have the same row count", rraIdx))
		}

		for rowIdx, row := range rra.Database.Rows {
			for vIdx, v := range row.V {
				if !floatEquals(v, xmlA.Rra[rraIdx].Database.Rows[rowIdx].V[vIdx]) {
					return errors.New(fmt.Sprintf("Expected RRA %d row %d to have the same values: %s != %s", rraIdx, rowIdx, v, xmlA.Rra[rraIdx].Database.Rows[rowIdx].V[vIdx]))
				}
			}
		}
	}

	return nil
}

func floatEquals(a, b string) bool {
	aF, aErr := strconv.ParseFloat(a, 64)
	bF, bErr := strconv.ParseFloat(b, 64)

	// 1e-3 is perfectly fine because our values are several orders of magnitude bigger
	return aErr == nil && bErr == nil && ((math.IsNaN(aF) && math.IsNaN(bF)) || (math.IsInf(aF, 0) && math.IsInf(bF, 0)) || math.Abs(bF-aF) < 1e-3)
}
