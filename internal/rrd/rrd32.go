//go:build rrd32
// +build rrd32

// This is a generated file! Please edit source .ksy file and use kaitai-struct-compiler to rebuild

package rrd

import (
	"bytes"
	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

type Rrd struct {
	Header       *Rrd_Header
	DsStore      []*Rrd_DsDef
	RraStore     []*Rrd_RraDef
	LiveHead     *Rrd_LiveHead
	PdpPrepStore []*Rrd_PdpPrep
	CdpPrepStore []interface{}
	RraPtrStore  []uint32
	RraDataStore []*Rrd_RraData
	_io          *kaitai.Stream
	_root        *Rrd
	_parent      interface{}
}

func NewRrd() *Rrd {
	return &Rrd{}
}

func (this *Rrd) Read(io *kaitai.Stream, parent interface{}, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp1 := NewRrd_Header()
	err = tmp1.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Header = tmp1
	this.DsStore = make([]*Rrd_DsDef, this.Header.DsCount)
	for i := range this.DsStore {
		tmp2 := NewRrd_DsDef()
		err = tmp2.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.DsStore[i] = tmp2
	}
	this.RraStore = make([]*Rrd_RraDef, this.Header.RraCount)
	for i := range this.RraStore {
		tmp3 := NewRrd_RraDef()
		err = tmp3.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.RraStore[i] = tmp3
	}
	tmp4 := NewRrd_LiveHead()
	err = tmp4.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.LiveHead = tmp4
	this.PdpPrepStore = make([]*Rrd_PdpPrep, this.Header.DsCount)
	for i := range this.PdpPrepStore {
		tmp5 := NewRrd_PdpPrep()
		err = tmp5.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.PdpPrepStore[i] = tmp5
	}
	this.CdpPrepStore = make([]interface{}, this.Header.RraCount)
	for i := range this.CdpPrepStore {
		switch this.RraStore[i].Cf {
		case "DEVPREDICT":
			tmp6 := NewRrd_UnknownCdpPrep()
			err = tmp6.Read(this._io, this, this._root)
			if err != nil {
				return err
			}
			this.CdpPrepStore[i] = tmp6
		case "FAILURES":
			tmp7 := NewRrd_UnknownCdpPrep()
			err = tmp7.Read(this._io, this, this._root)
			if err != nil {
				return err
			}
			this.CdpPrepStore[i] = tmp7
		case "SEASONAL":
			tmp8 := NewRrd_UnknownCdpPrep()
			err = tmp8.Read(this._io, this, this._root)
			if err != nil {
				return err
			}
			this.CdpPrepStore[i] = tmp8
		case "DEVSEASONAL":
			tmp9 := NewRrd_UnknownCdpPrep()
			err = tmp9.Read(this._io, this, this._root)
			if err != nil {
				return err
			}
			this.CdpPrepStore[i] = tmp9
		case "MHWPREDICT":
			tmp10 := NewRrd_UnknownCdpPrep()
			err = tmp10.Read(this._io, this, this._root)
			if err != nil {
				return err
			}
			this.CdpPrepStore[i] = tmp10
		case "HWPREDICT":
			tmp11 := NewRrd_UnknownCdpPrep()
			err = tmp11.Read(this._io, this, this._root)
			if err != nil {
				return err
			}
			this.CdpPrepStore[i] = tmp11
		default:
			tmp12 := NewRrd_CdpPrep()
			err = tmp12.Read(this._io, this, this._root)
			if err != nil {
				return err
			}
			this.CdpPrepStore[i] = tmp12
		}
	}
	this.RraPtrStore = make([]uint32, this.Header.RraCount)
	for i := range this.RraPtrStore {
		tmp13, err := this._io.ReadU4le()
		if err != nil {
			return err
		}
		this.RraPtrStore[i] = tmp13
	}
	this.RraDataStore = make([]*Rrd_RraData, this.Header.RraCount)
	for i := range this.RraDataStore {
		tmp14 := NewRrd_RraData(this.RraStore[i].RowCount)
		err = tmp14.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.RraDataStore[i] = tmp14
	}
	return err
}

type Rrd_PdpPrepParams struct {
	UnknownSecCount uint32
	_unnamed1       []byte
	CurrentValue    float64
	Unused          [][]byte
	_io             *kaitai.Stream
	_root           *Rrd
	_parent         *Rrd_PdpPrep
}

func NewRrd_PdpPrepParams() *Rrd_PdpPrepParams {
	return &Rrd_PdpPrepParams{}
}

func (this *Rrd_PdpPrepParams) Read(io *kaitai.Stream, parent *Rrd_PdpPrep, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp15, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.UnknownSecCount = uint32(tmp15)
	tmp16, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp16 = tmp16
	this._unnamed1 = tmp16
	tmp17, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.CurrentValue = float64(tmp17)
	this.Unused = make([][]byte, (10 - 2))
	for i := range this.Unused {
		tmp18, err := this._io.ReadBytes(int(8))
		if err != nil {
			return err
		}
		tmp18 = tmp18
		this.Unused[i] = tmp18
	}
	return err
}

/**
 * padding
 */
type Rrd_RraParams struct {
	Xff     float64
	Unused  [][]byte
	_io     *kaitai.Stream
	_root   *Rrd
	_parent *Rrd_RraDef
}

func NewRrd_RraParams() *Rrd_RraParams {
	return &Rrd_RraParams{}
}

func (this *Rrd_RraParams) Read(io *kaitai.Stream, parent *Rrd_RraDef, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp19, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.Xff = float64(tmp19)
	this.Unused = make([][]byte, (10 - 1))
	for i := range this.Unused {
		tmp20, err := this._io.ReadBytes(int(8))
		if err != nil {
			return err
		}
		tmp20 = tmp20
		this.Unused[i] = tmp20
	}
	return err
}

type Rrd_CdpPrepParams struct {
	Value           float64
	UnknownPdpCount uint32
	Unused          [][]byte
	_unnamed3       []byte
	PrimaryValue    float64
	SecondaryValue  float64
	_io             *kaitai.Stream
	_root           *Rrd
	_parent         *Rrd_CdpPrep
}

func NewRrd_CdpPrepParams() *Rrd_CdpPrepParams {
	return &Rrd_CdpPrepParams{}
}

func (this *Rrd_CdpPrepParams) Read(io *kaitai.Stream, parent *Rrd_CdpPrep, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp21, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.Value = float64(tmp21)
	tmp22, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.UnknownPdpCount = uint32(tmp22)
	this.Unused = make([][]byte, (10 - 4))
	for i := range this.Unused {
		tmp23, err := this._io.ReadBytes(int(8))
		if err != nil {
			return err
		}
		tmp23 = tmp23
		this.Unused[i] = tmp23
	}
	tmp24, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp24 = tmp24
	this._unnamed3 = tmp24
	tmp25, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.PrimaryValue = float64(tmp25)
	tmp26, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.SecondaryValue = float64(tmp26)
	return err
}

/**
 * padding
 */
type Rrd_RraDataRow struct {
	Values  []float64
	_io     *kaitai.Stream
	_root   *Rrd
	_parent *Rrd_RraData
}

func NewRrd_RraDataRow() *Rrd_RraDataRow {
	return &Rrd_RraDataRow{}
}

func (this *Rrd_RraDataRow) Read(io *kaitai.Stream, parent *Rrd_RraData, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	this.Values = make([]float64, this._root.Header.DsCount)
	for i := range this.Values {
		tmp27, err := this._io.ReadF8le()
		if err != nil {
			return err
		}
		this.Values[i] = tmp27
	}
	return err
}

type Rrd_LiveHead struct {
	LastUpdate     uint32
	LastUpdateUsec uint32
	_io            *kaitai.Stream
	_root          *Rrd
	_parent        *Rrd
}

func NewRrd_LiveHead() *Rrd_LiveHead {
	return &Rrd_LiveHead{}
}

func (this *Rrd_LiveHead) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp28, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.LastUpdate = uint32(tmp28)
	tmp29, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.LastUpdateUsec = uint32(tmp29)
	return err
}

type Rrd_UnknownParams struct {
	Params  [][]byte
	_io     *kaitai.Stream
	_root   *Rrd
	_parent *Rrd_RraDef
}

func NewRrd_UnknownParams() *Rrd_UnknownParams {
	return &Rrd_UnknownParams{}
}

func (this *Rrd_UnknownParams) Read(io *kaitai.Stream, parent *Rrd_RraDef, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	this.Params = make([][]byte, 10)
	for i := range this.Params {
		tmp30, err := this._io.ReadBytes(int(8))
		if err != nil {
			return err
		}
		tmp30 = tmp30
		this.Params[i] = tmp30
	}
	return err
}

type Rrd_PdpPrep struct {
	LastDsReading string
	Params        *Rrd_PdpPrepParams
	_io           *kaitai.Stream
	_root         *Rrd
	_parent       *Rrd
}

func NewRrd_PdpPrep() *Rrd_PdpPrep {
	return &Rrd_PdpPrep{}
}

func (this *Rrd_PdpPrep) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp31, err := this._io.ReadBytes(int(32))
	if err != nil {
		return err
	}
	tmp31 = kaitai.BytesTerminate(tmp31, 0, false)
	this.LastDsReading = string(tmp31)
	tmp32 := NewRrd_PdpPrepParams()
	err = tmp32.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Params = tmp32
	return err
}

type Rrd_CdpPrep struct {
	Params  []*Rrd_CdpPrepParams
	_io     *kaitai.Stream
	_root   *Rrd
	_parent *Rrd
}

func NewRrd_CdpPrep() *Rrd_CdpPrep {
	return &Rrd_CdpPrep{}
}

func (this *Rrd_CdpPrep) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	this.Params = make([]*Rrd_CdpPrepParams, this._root.Header.DsCount)
	for i := range this.Params {
		tmp33 := NewRrd_CdpPrepParams()
		err = tmp33.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params[i] = tmp33
	}
	return err
}

type Rrd_DsDef struct {
	Name       string
	Datasource string
	Params     *Rrd_DsParams
	_io        *kaitai.Stream
	_root      *Rrd
	_parent    *Rrd
}

func NewRrd_DsDef() *Rrd_DsDef {
	return &Rrd_DsDef{}
}

func (this *Rrd_DsDef) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp34, err := this._io.ReadBytes(int(20))
	if err != nil {
		return err
	}
	tmp34 = kaitai.BytesTerminate(tmp34, 0, false)
	this.Name = string(tmp34)
	tmp35, err := this._io.ReadBytes(int(20))
	if err != nil {
		return err
	}
	tmp35 = kaitai.BytesTerminate(tmp35, 0, false)
	this.Datasource = string(tmp35)
	tmp36 := NewRrd_DsParams()
	err = tmp36.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Params = tmp36
	return err
}

type Rrd_Header struct {
	Magic       []byte
	Version     string
	Padding     []byte
	Floatcookie float64
	DsCount     uint32
	RraCount    uint32
	PdpStep     uint32
	_unnamed7   []byte
	Params      [][]byte
	_io         *kaitai.Stream
	_root       *Rrd
	_parent     *Rrd
}

func NewRrd_Header() *Rrd_Header {
	return &Rrd_Header{}
}

func (this *Rrd_Header) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp37, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp37 = tmp37
	this.Magic = tmp37
	if !(bytes.Equal(this.Magic, []uint8{82, 82, 68, 0})) {
		return kaitai.NewValidationNotEqualError([]uint8{82, 82, 68, 0}, this.Magic, this._io, "/types/header/seq/0")
	}
	tmp38, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp38 = tmp38
	this.Version = string(tmp38)
	tmp39, err := this._io.ReadBytes(int(8))
	if err != nil {
		return err
	}
	tmp39 = tmp39
	this.Padding = tmp39
	tmp40, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.Floatcookie = float64(tmp40)
	tmp41, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.DsCount = uint32(tmp41)
	tmp42, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.RraCount = uint32(tmp42)
	tmp43, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.PdpStep = uint32(tmp43)
	tmp44, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp44 = tmp44
	this._unnamed7 = tmp44
	this.Params = make([][]byte, 10)
	for i := range this.Params {
		tmp45, err := this._io.ReadBytes(int(8))
		if err != nil {
			return err
		}
		tmp45 = tmp45
		this.Params[i] = tmp45
	}
	return err
}

/**
 * padding
 */

/**
 * unused
 */
type Rrd_DsParams struct {
	MinHeartbeatCount uint32
	MinVal            float64
	MaxVal            float64
	Cdef              uint32
	Unused            [][]byte
	_io               *kaitai.Stream
	_root             *Rrd
	_parent           *Rrd_DsDef
}

func NewRrd_DsParams() *Rrd_DsParams {
	return &Rrd_DsParams{}
}

func (this *Rrd_DsParams) Read(io *kaitai.Stream, parent *Rrd_DsDef, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp46, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.MinHeartbeatCount = uint32(tmp46)
	tmp47, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.MinVal = float64(tmp47)
	tmp48, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.MaxVal = float64(tmp48)
	tmp49, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.Cdef = uint32(tmp49)
	this.Unused = make([][]byte, (10 - 3))
	for i := range this.Unused {
		tmp50, err := this._io.ReadBytes(int(8))
		if err != nil {
			return err
		}
		tmp50 = tmp50
		this.Unused[i] = tmp50
	}
	return err
}

/**
 * pointer to encoded rpn expression only applies to DST_CDEF, not supported here
 */
type Rrd_RraData struct {
	Row      []*Rrd_RraDataRow
	RowCount uint32
	_io      *kaitai.Stream
	_root    *Rrd
	_parent  *Rrd
}

func NewRrd_RraData(rowCount uint32) *Rrd_RraData {
	return &Rrd_RraData{
		RowCount: rowCount,
	}
}

func (this *Rrd_RraData) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	this.Row = make([]*Rrd_RraDataRow, this.RowCount)
	for i := range this.Row {
		tmp51 := NewRrd_RraDataRow()
		err = tmp51.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Row[i] = tmp51
	}
	return err
}

type Rrd_RraDef struct {
	Cf        string
	RowCount  uint32
	PdpCount  uint32
	_unnamed3 []byte
	Params    interface{}
	_io       *kaitai.Stream
	_root     *Rrd
	_parent   *Rrd
}

func NewRrd_RraDef() *Rrd_RraDef {
	return &Rrd_RraDef{}
}

func (this *Rrd_RraDef) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp52, err := this._io.ReadBytes(int(20))
	if err != nil {
		return err
	}
	tmp52 = kaitai.BytesTerminate(tmp52, 0, false)
	this.Cf = string(tmp52)
	tmp53, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.RowCount = uint32(tmp53)
	tmp54, err := this._io.ReadU4le()
	if err != nil {
		return err
	}
	this.PdpCount = uint32(tmp54)
	tmp55, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp55 = tmp55
	this._unnamed3 = tmp55
	switch this.Cf {
	case "DEVPREDICT":
		tmp56 := NewRrd_UnknownParams()
		err = tmp56.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp56
	case "FAILURES":
		tmp57 := NewRrd_UnknownParams()
		err = tmp57.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp57
	case "SEASONAL":
		tmp58 := NewRrd_UnknownParams()
		err = tmp58.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp58
	case "DEVSEASONAL":
		tmp59 := NewRrd_UnknownParams()
		err = tmp59.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp59
	case "MHWPREDICT":
		tmp60 := NewRrd_UnknownParams()
		err = tmp60.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp60
	case "HWPREDICT":
		tmp61 := NewRrd_UnknownParams()
		err = tmp61.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp61
	default:
		tmp62 := NewRrd_RraParams()
		err = tmp62.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp62
	}
	return err
}

/**
 * padding
 */
type Rrd_UnknownCdpPrep struct {
	Unused  [][]byte
	_io     *kaitai.Stream
	_root   *Rrd
	_parent *Rrd
}

func NewRrd_UnknownCdpPrep() *Rrd_UnknownCdpPrep {
	return &Rrd_UnknownCdpPrep{}
}

func (this *Rrd_UnknownCdpPrep) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	this.Unused = make([][]byte, 10)
	for i := range this.Unused {
		tmp63, err := this._io.ReadBytes(int(8))
		if err != nil {
			return err
		}
		tmp63 = tmp63
		this.Unused[i] = tmp63
	}
	return err
}
