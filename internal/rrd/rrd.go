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
	RraPtrStore  []uint64
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
	this.RraPtrStore = make([]uint64, this.Header.RraCount)
	for i := range this.RraPtrStore {
		tmp13, err := this._io.ReadU8le()
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
	UnknownSecCount uint64
	CurrentValue    float64
	Unused          []uint64
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

	tmp15, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.UnknownSecCount = uint64(tmp15)
	tmp16, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.CurrentValue = float64(tmp16)
	this.Unused = make([]uint64, (10 - 2))
	for i := range this.Unused {
		tmp17, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.Unused[i] = tmp17
	}
	return err
}

type Rrd_RraParams struct {
	Xff     float64
	Unused  []uint64
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

	tmp18, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.Xff = float64(tmp18)
	this.Unused = make([]uint64, (10 - 1))
	for i := range this.Unused {
		tmp19, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.Unused[i] = tmp19
	}
	return err
}

type Rrd_CdpPrepParams struct {
	Value           float64
	UnknownPdpCount uint64
	Unused          []uint64
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

	tmp20, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.Value = float64(tmp20)
	tmp21, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.UnknownPdpCount = uint64(tmp21)
	this.Unused = make([]uint64, (10 - 4))
	for i := range this.Unused {
		tmp22, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.Unused[i] = tmp22
	}
	tmp23, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.PrimaryValue = float64(tmp23)
	tmp24, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.SecondaryValue = float64(tmp24)
	return err
}

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
		tmp25, err := this._io.ReadF8le()
		if err != nil {
			return err
		}
		this.Values[i] = tmp25
	}
	return err
}

type Rrd_LiveHead struct {
	LastUpdate     uint64
	LastUpdateUsec uint64
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

	tmp26, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.LastUpdate = uint64(tmp26)
	tmp27, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.LastUpdateUsec = uint64(tmp27)
	return err
}

type Rrd_UnknownParams struct {
	Params  []uint64
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

	this.Params = make([]uint64, 10)
	for i := range this.Params {
		tmp28, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.Params[i] = tmp28
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

	tmp29, err := this._io.ReadBytes(int(32))
	if err != nil {
		return err
	}
	tmp29 = kaitai.BytesTerminate(tmp29, 0, false)
	this.LastDsReading = string(tmp29)
	tmp30 := NewRrd_PdpPrepParams()
	err = tmp30.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Params = tmp30
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
		tmp31 := NewRrd_CdpPrepParams()
		err = tmp31.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params[i] = tmp31
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

	tmp32, err := this._io.ReadBytes(int(20))
	if err != nil {
		return err
	}
	tmp32 = kaitai.BytesTerminate(tmp32, 0, false)
	this.Name = string(tmp32)
	tmp33, err := this._io.ReadBytes(int(20))
	if err != nil {
		return err
	}
	tmp33 = kaitai.BytesTerminate(tmp33, 0, false)
	this.Datasource = string(tmp33)
	tmp34 := NewRrd_DsParams()
	err = tmp34.Read(this._io, this, this._root)
	if err != nil {
		return err
	}
	this.Params = tmp34
	return err
}

type Rrd_Header struct {
	Magic       []byte
	Version     string
	Padding     []byte
	Floatcookie float64
	DsCount     uint64
	RraCount    uint64
	PdpStep     uint64
	Params      []uint64
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

	tmp35, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp35 = tmp35
	this.Magic = tmp35
	if !(bytes.Equal(this.Magic, []uint8{82, 82, 68, 0})) {
		return kaitai.NewValidationNotEqualError([]uint8{82, 82, 68, 0}, this.Magic, this._io, "/types/header/seq/0")
	}
	tmp36, err := this._io.ReadBytes(int(4))
	if err != nil {
		return err
	}
	tmp36 = tmp36
	this.Version = string(tmp36)
	tmp37, err := this._io.ReadBytes(int(8))
	if err != nil {
		return err
	}
	tmp37 = tmp37
	this.Padding = tmp37
	tmp38, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.Floatcookie = float64(tmp38)
	tmp39, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.DsCount = uint64(tmp39)
	tmp40, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.RraCount = uint64(tmp40)
	tmp41, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.PdpStep = uint64(tmp41)
	this.Params = make([]uint64, 10)
	for i := range this.Params {
		tmp42, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.Params[i] = tmp42
	}
	return err
}

/**
 * unused
 */
type Rrd_DsParams struct {
	MinHeartbeatCount uint64
	MinVal            float64
	MaxVal            float64
	Cdef              uint64
	Unused            []uint64
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

	tmp43, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.MinHeartbeatCount = uint64(tmp43)
	tmp44, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.MinVal = float64(tmp44)
	tmp45, err := this._io.ReadF8le()
	if err != nil {
		return err
	}
	this.MaxVal = float64(tmp45)
	tmp46, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.Cdef = uint64(tmp46)
	this.Unused = make([]uint64, (10 - 4))
	for i := range this.Unused {
		tmp47, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.Unused[i] = tmp47
	}
	return err
}

/**
 * pointer to encoded rpn expression only applies to DST_CDEF, not supported here
 */
type Rrd_RraData struct {
	Row      []*Rrd_RraDataRow
	RowCount uint64
	_io      *kaitai.Stream
	_root    *Rrd
	_parent  *Rrd
}

func NewRrd_RraData(rowCount uint64) *Rrd_RraData {
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
		tmp48 := NewRrd_RraDataRow()
		err = tmp48.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Row[i] = tmp48
	}
	return err
}

type Rrd_RraDef struct {
	Cf       string
	RowCount uint64
	PdpCount uint64
	Params   interface{}
	_io      *kaitai.Stream
	_root    *Rrd
	_parent  *Rrd
}

func NewRrd_RraDef() *Rrd_RraDef {
	return &Rrd_RraDef{}
}

func (this *Rrd_RraDef) Read(io *kaitai.Stream, parent *Rrd, root *Rrd) (err error) {
	this._io = io
	this._parent = parent
	this._root = root

	tmp49, err := this._io.ReadBytes(int(24))
	if err != nil {
		return err
	}
	tmp49 = kaitai.BytesTerminate(tmp49, 0, false)
	this.Cf = string(tmp49)
	tmp50, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.RowCount = uint64(tmp50)
	tmp51, err := this._io.ReadU8le()
	if err != nil {
		return err
	}
	this.PdpCount = uint64(tmp51)
	switch this.Cf {
	case "DEVPREDICT":
		tmp52 := NewRrd_UnknownParams()
		err = tmp52.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp52
	case "FAILURES":
		tmp53 := NewRrd_UnknownParams()
		err = tmp53.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp53
	case "SEASONAL":
		tmp54 := NewRrd_UnknownParams()
		err = tmp54.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp54
	case "DEVSEASONAL":
		tmp55 := NewRrd_UnknownParams()
		err = tmp55.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp55
	case "MHWPREDICT":
		tmp56 := NewRrd_UnknownParams()
		err = tmp56.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp56
	case "HWPREDICT":
		tmp57 := NewRrd_UnknownParams()
		err = tmp57.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp57
	default:
		tmp58 := NewRrd_RraParams()
		err = tmp58.Read(this._io, this, this._root)
		if err != nil {
			return err
		}
		this.Params = tmp58
	}
	return err
}

type Rrd_UnknownCdpPrep struct {
	Unused  []uint64
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

	this.Unused = make([]uint64, 10)
	for i := range this.Unused {
		tmp59, err := this._io.ReadU8le()
		if err != nil {
			return err
		}
		this.Unused[i] = tmp59
	}
	return err
}
