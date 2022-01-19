//go:build rrdtool
// +build rrdtool

package rrdtool

import (
	"encoding/xml"
	"fmt"
	"io"
	"os/exec"

	"github.com/scaleway/rrdmerge/internal/rrd"
)

func Dump(file string) (*rrd.Rrd, error) {
	rrdtoolCmd := exec.Command("rrdtool", "dump", file)
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

func Flush(file string, daemonOpt string) error {
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
			return fmt.Errorf("rrdtool returned: %s", out)
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

func Restore(rrdPtr *rrd.Rrd, file string) error {
	xmlBytes, err := xml.Marshal(rrd.FromRRDStruct(rrdPtr))
	if err != nil {
		return err
	}

	rrdtoolCmd := exec.Command("rrdtool", "restore", "-f", "-", file)
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

	if out, err := io.ReadAll(stdErr); err == nil {
		if len(out) > 0 {
			return fmt.Errorf("rrdtool returned: %s", out)
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
