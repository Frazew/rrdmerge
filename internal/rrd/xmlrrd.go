package rrd

import (
	"encoding/xml"
	"fmt"
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
