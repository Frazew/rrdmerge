//go:build librrd
// +build librrd

package rrdtool

// #cgo LDFLAGS: -lrrd
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include <rrd.h>
import "C"

import (
	"encoding/xml"
	"errors"
	"os"
	"unsafe"

	"github.com/scaleway/rrdmerge/internal/rrd"
)

func Flush(filename string, daemonOpt string) (err error) {
	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	argv := []string{
		"-d",
		daemonOpt,
		filename,
	}
	cArgv := makeCStringArray(argv)
	defer freeCStringArray(cArgv)

	C.rrd_clear_error()
	ret := C.rrd_flushcached(C.int(len(argv)), getCStringArrayPointer(cArgv))

	if int(ret) != 0 {
		err = errors.New(C.GoString(C.rrd_get_error()))
	}
	return err
}

func Restore(rrdPtr *rrd.Rrd, filename string) (err error) {
	xmlBytes, err := xml.Marshal(rrd.FromRRDStruct(rrdPtr))
	if err != nil {
		return err
	}
	file, err := os.CreateTemp("/tmp/", "rrdmerge")
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()
	if err != nil {
		return err
	}

	_, err = file.Write(xmlBytes)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}

	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	argv := []string{
		"--range-check",
		"--force-overwrite",
		file.Name(),
		filename,
	}
	cArgv := makeCStringArray(argv)
	defer freeCStringArray(cArgv)

	C.rrd_clear_error()
	ret := C.rrd_restore(C.int(len(argv)), getCStringArrayPointer(cArgv))

	if int(ret) != 0 {
		err = errors.New(C.GoString(C.rrd_get_error()))
	}
	return err
}

func makeCStringArray(values []string) (cvalues []*C.char) {
	cvalues = make([]*C.char, len(values))
	for i := range values {
		cvalues[i] = C.CString(values[i])
	}
	return
}

func freeCStringArray(cvalues []*C.char) {
	for i := range cvalues {
		C.free(unsafe.Pointer(cvalues[i]))
	}
}

func getCStringArrayPointer(cvalues []*C.char) **C.char {
	return (**C.char)(unsafe.Pointer(&cvalues[0]))
}
