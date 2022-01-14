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

func Load(file string) (*rrd.Rrd, error) {
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

func Dump(rrd *rrd.Rrd, file string, mode os.FileMode) error {
	xmlBytes, err := xml.Marshal(rrd)
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
