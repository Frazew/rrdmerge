package rrdtool

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/scaleway/rrdmerge/internal/rrd"
)

func Dump(file string) (*rrd.Rrd, error) {
	rrdtoolCmd := exec.Command("rrdtool", "dump", file, "/dev/stdout")
	stdOut, err := rrdtoolCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = rrdtoolCmd.Start()
	if err != nil {
		return nil, err
	}

	var rrd rrd.Rrd
	if err := xml.NewDecoder(stdOut).Decode(&rrd); err != nil {
		return nil, err
	}

	err = rrdtoolCmd.Wait()
	if err != nil {
		return nil, err
	}

	return &rrd, nil
}

func Tune(file string, daemonOpt string) error {
	rrdtoolCmd := exec.Command("rrdtool", "flushcached", file, "-d", daemonOpt)
	stdErr, err := rrdtoolCmd.StderrPipe()
	if err != nil {
		return err
	}

	err = rrdtoolCmd.Start()
	if err != nil {
		return err
	}

	if out, err := io.ReadAll(stdErr); err == nil {
		if len(out) > 0 {
			return errors.New(fmt.Sprintf("rrdtool returned: %s", out))
		}
	} else {
		return err
	}

	err = rrdtoolCmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func Restore(rrdPtr *rrd.Rrd, file string, mode os.FileMode) error {
	dsSlice := make([]rrd.XmlDs, rrdPtr.Header.DsCount)
	for dsIdx, ds := range rrdPtr.DsStore {
		dsSlice[dsIdx] = rrd.XmlDs{
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

	rraSlice := make([]rrd.XmlRra, rrdPtr.Header.RraCount)
	for rraIdx, rra := range rrdPtr.RraDataStore {
		dsParams := make([]rrd.XmlDsParams, rrdPtr.Header.DsCount)
		for cdpIdx, cdp := range rrdPtr.CdpPrepStore[rraIdx].(*rrd.Rrd_CdpPrep).Params {
			dsParams[cdpIdx] = rrd.XmlDsParams{
				PrimaryValue:      fmt.Sprintf("%.10e", cdp.PrimaryValue),
				SecondaryValue:    fmt.Sprintf("%.10e", cdp.SecondaryValue),
				Value:             fmt.Sprintf("%.10e", cdp.Value),
				UnknownDatapoints: fmt.Sprintf("%d", cdp.UnknownPdpCount),
			}
		}

		rows := make([]struct {
			V []string "xml:\"v\""
		}, rra.RowCount)

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

		rraSlice[rraIdx] = rrd.XmlRra{
			Cf:        rrdPtr.RraStore[rraIdx].Cf,
			PdpPerRow: int(rrdPtr.RraStore[rraIdx].PdpCount),
			Params: struct {
				Xff string "xml:\"xff\""
			}{fmt.Sprintf("%.10e", rrdPtr.RraStore[rraIdx].Params.(*rrd.Rrd_RraParams).Xff)},
			CdpPrep: struct {
				Ds []rrd.XmlDsParams "xml:\"ds\""
			}{Ds: dsParams},
			Database: struct {
				Row []struct {
					V []string "xml:\"v\""
				} "xml:\"row\""
			}{Row: rows},
		}
	}

	xmlRrd := rrd.XmlRrd{
		Version:    "0003",
		Step:       int(rrdPtr.Header.PdpStep),
		Lastupdate: int64(rrdPtr.LiveHead.LastUpdate),
		Ds:         dsSlice,
		Rra:        rraSlice,
	}
	xmlBytes, err := xml.Marshal(xmlRrd)
	if err != nil {
		return err
	}

	rrdtoolCmd := exec.Command("rrdtool", "restore", "-f", "/dev/stdin", file)
	stdIn, err := rrdtoolCmd.StdinPipe()
	if err != nil {
		return err
	}
	stdErr, err := rrdtoolCmd.StderrPipe()
	if err != nil {
		return err
	}

	err = rrdtoolCmd.Start()
	if err != nil {
		return err
	}

	go func() {
		_, err = stdIn.Write(xmlBytes)
		err = stdIn.Close()
	}()

	f, err := os.Create("/tmp/out/" + filepath.Base(file))
	f.Write(xmlBytes)
	f.Sync()
	f.Close()

	if out, err := io.ReadAll(stdErr); err == nil {
		if len(out) > 0 {
			return errors.New(fmt.Sprintf("rrdtool returned: %s", out))
		}
	} else {
		return err
	}

	rrdtoolCmd.Wait()
	if err != nil {
		return err
	}

	return nil
}
