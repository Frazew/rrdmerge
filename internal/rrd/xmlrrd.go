package rrd

import "encoding/xml"

type XmlDsParams struct {
	PrimaryValue      string `xml:"primary_value"`
	SecondaryValue    string `xml:"secondary_value"`
	Value             string `xml:"value"`
	UnknownDatapoints string `xml:"unknown_datapoints"`
}

type XmlRra struct {
	Cf        string   `xml:"cf"`
	PdpPerRow int      `xml:"pdp_per_row"`
	Params    struct { // Only support CF_AVERAGE, CF_MAXIMUM, CF_MINIMUM, and CF_LAST
		Xff string `xml:"xff"`
	} `xml:"params"`
	CdpPrep struct {
		Ds []XmlDsParams `xml:"ds"`
	} `xml:"cdp_prep"`
	Database struct {
		Row []struct {
			V []string `xml:"v"`
		} `xml:"row"`
	} `xml:"database"`
}

type XmlDs struct {
	Name             string `xml:"name"`
	Type             string `xml:"type"`
	MinimalHeartbeat string `xml:"minimal_heartbeat"`
	Min              string `xml:"min"`
	Max              string `xml:"max"`
	LastDs           string `xml:"last_ds"`
	Value            string `xml:"value"`
	UnknownSec       string `xml:"unknown_sec"`
}

type XmlRrd struct {
	XMLName    xml.Name `xml:"rrd"`
	Version    string   `xml:"version"`
	Step       int      `xml:"step"`
	Lastupdate int64    `xml:"lastupdate"`
	Ds         []XmlDs  `xml:"ds"`
	Rra        []XmlRra `xml:"rra"`
}
