package binchunk

import (
	"encoding/binary"
	"math"
)

// Reader binary chunk reader
type Reader struct {
	data []byte
}

// ReadByte read byte
func (imp *Reader) ReadByte() byte {
	b := imp.data[0]
	imp.data = imp.data[1:]
	return b
}

// ReadBytes read bytes
func (imp *Reader) ReadBytes(n uint) []byte {
	bytes := imp.data[:n]
	imp.data = imp.data[n:]
	return bytes
}

// ReadUint32 read uint32
func (imp *Reader) ReadUint32() uint32 {
	n := binary.LittleEndian.Uint32(imp.data)
	imp.data = imp.data[4:]
	return n
}

// ReadUint64 read uint64
func (imp *Reader) ReadUint64() uint64 {
	n := binary.LittleEndian.Uint64(imp.data)
	imp.data = imp.data[8:]
	return n
}

// ReadLuaInteger read lua integer
func (imp *Reader) ReadLuaInteger() int64 {
	return int64(imp.ReadUint64())
}

// ReadLuaNumber read lua number
func (imp *Reader) ReadLuaNumber() float64 {
	return math.Float64frombits(imp.ReadUint64())
}

// ReadString read string
func (imp *Reader) ReadString() string {
	size := uint(imp.ReadByte())
	if size <= 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(imp.ReadUint64())
	}
	bytes := imp.ReadBytes(size - 1)
	return string(bytes)
}

// CheckHeader check header
func (imp *Reader) CheckHeader() {
	if string(imp.ReadBytes(4)) != LuaSignature {
		panic("not a precompiled chunk!")
	}
	if imp.ReadByte() != LuacVersion {
		panic("version mismatch!")
	}
	if imp.ReadByte() != LuacFormat {
		panic("format mismatch!")
	}
	if string(imp.ReadBytes(6)) != LuacData {
		panic("corrupted!")
	}
	if imp.ReadByte() != CIntSize {
		panic("int size mismatch!")
	}
	if imp.ReadByte() != CSizetSize {
		panic("size_t size mismatch!")
	}
	if imp.ReadByte() != InstructionSize {
		panic("instruction size mismatch!")
	}
	if imp.ReadByte() != LuaIntSize {
		panic("lua_Integer size mismatch!")
	}
	if imp.ReadByte() != LuaNumberSize {
		panic("lua_Number size mismatch!")
	}
	if imp.ReadLuaInteger() != LuacInt {
		panic("endianness mismatch!")
	}
	if imp.ReadLuaNumber() != LuacNum {
		panic("float format mismatch!")
	}
}

// ReadProto read proto
func (imp *Reader) ReadProto(parentSource string) *ProtoType {
	source := imp.ReadString()
	if source == "" {
		source = parentSource
	}
	return &ProtoType{
		Source:          source,
		LineDefined:     imp.ReadUint32(),
		LastLineDefined: imp.ReadUint32(),
		NumParams:       imp.ReadByte(),
		IsVararg:        imp.ReadByte(),
		MaxStackSize:    imp.ReadByte(),
		Code:            imp.ReadCode(),
		Constants:       imp.ReadConstants(),
		Upvalues:        imp.ReadUpvalues(),
		Protos:          imp.ReadProtos(source),
		LineInfo:        imp.ReadLineInfo(),
		LocVars:         imp.ReadLocVars(),
		UpvalueNames:    imp.ReadUpvalueNames(),
	}
}

// ReadCode read code
func (imp *Reader) ReadCode() []uint32 {
	code := make([]uint32, imp.ReadUint32())
	for i := range code {
		code[i] = imp.ReadUint32()
	}
	return code
}

// ReadConstants read constants
func (imp *Reader) ReadConstants() []interface{} {
	constants := make([]interface{}, imp.ReadUint32())
	for i := range constants {
		constants[i] = imp.ReadConstant()
	}
	return constants
}

// ReadConstant read constant
func (imp *Reader) ReadConstant() interface{} {
	switch imp.ReadByte() {
	case TagNil:
		return nil
	case TagBoolean:
		return imp.ReadByte() != 0
	case TagInteger:
		return imp.ReadLuaInteger()
	case TagNumber:
		return imp.ReadLuaNumber()
	case TagShortStr, TagLongStr:
		return imp.ReadString()
	default:
		panic("corrupted") // todo
	}
}

// ReadUpvalues read upvalues
func (imp *Reader) ReadUpvalues() []Upvalue {
	upvalues := make([]Upvalue, imp.ReadUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: imp.ReadByte(),
			Idx:     imp.ReadByte(),
		}
	}
	return upvalues
}

// ReadProtos read protos
func (imp *Reader) ReadProtos(parentSource string) []*ProtoType {
	protos := make([]*ProtoType, imp.ReadUint32())
	for i := range protos {
		protos[i] = imp.ReadProto(parentSource)
	}
	return protos
}

// ReadLineInfo read line info
func (imp *Reader) ReadLineInfo() []uint32 {
	lineInfo := make([]uint32, imp.ReadUint32())
	for i := range lineInfo {
		lineInfo[i] = imp.ReadUint32()
	}
	return lineInfo
}

// ReadLocVars read local vars
func (imp *Reader) ReadLocVars() []LocVar {
	locVars := make([]LocVar, imp.ReadUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: imp.ReadString(),
			StartPC: imp.ReadUint32(),
			EndPC:   imp.ReadUint32(),
		}
	}
	return locVars
}

// ReadUpvalueNames read upvalues name
func (imp *Reader) ReadUpvalueNames() []string {
	names := make([]string, imp.ReadUint32())
	for i := range names {
		names[i] = imp.ReadString()
	}
	return names
}
