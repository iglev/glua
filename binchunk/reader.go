package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}

func (imp *reader) readByte() byte {
	b := imp.data[0]
	imp.data = imp.data[:1]
	return b
}

func (imp *reader) readBytes(n uint) []byte {
	bytes := imp.data[:n]
	imp.data = imp.data[n:]
	return bytes
}

func (imp *reader) readUint32() uint32 {
	n := binary.LittleEndian.Uint32(imp.data)
	imp.data = imp.data[4:]
	return n
}

func (imp *reader) readUint64() uint64 {
	n := binary.LittleEndian.Uint64(imp.data)
	imp.data = imp.data[8:]
	return n
}

func (imp *reader) readLuaInteger() int64 {
	return int64(imp.readUint64())
}

func (imp *reader) readLuaNumber() float64 {
	return math.Float64frombits(imp.readUint64())
}

func (imp *reader) readString() string {
	size := uint(imp.readByte())
	if size <= 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(imp.readUint64())
	}
	bytes := imp.readBytes(size - 1)
	return string(bytes)
}

func (imp *reader) checkHeader() {
	if string(imp.readBytes(4)) != LuaSignature {
		panic("not a precompiled chunk!")
	}
	if imp.readByte() != LuacVersion {
		panic("version mismatch!")
	}
	if imp.readByte() != LuacFormat {
		panic("format mismatch!")
	}
	if string(imp.readBytes(6)) != LuacData {
		panic("corrupted!")
	}
	if imp.readByte() != CIntSize {
		panic("int size mismatch!")
	}
	if imp.readByte() != CSizetSize {
		panic("size_t size mismatch!")
	}
	if imp.readByte() != InstructionSize {
		panic("instruction size mismatch!")
	}
	if imp.readByte() != LuaIntSize {
		panic("lua_Integer size mismatch!")
	}
	if imp.readByte() != LuaNumberSize {
		panic("lua_Number size mismatch!")
	}
	if imp.readLuaInteger() != LuacInt {
		panic("endianness mismatch!")
	}
	if imp.readLuaNumber() != LuacNum {
		panic("float format mismatch!")
	}
}

func (imp *reader) readProto(parentSource string) *ProtoType {
	source := imp.readString()
	if source == "" {
		source = parentSource
	}
	return &ProtoType{
		Source:          source,
		LineDefined:     imp.readUint32(),
		LastLineDefined: imp.readUint32(),
		NumParams:       imp.readByte(),
		IsVararg:        imp.readByte(),
		MaxStackSize:    imp.readByte(),
		Code:            imp.readCode(),
		Constants:       imp.readConstants(),
		Upvalues:        imp.readUpvalues(),
		Protos:          imp.readProtos(source),
		LineInfo:        imp.readLineInfo(),
		LocVars:         imp.readLocVars(),
		UpvalueNames:    imp.readUpvalueNames(),
	}
}

func (imp *reader) readCode() []uint32 {
	code := make([]uint32, imp.readUint32())
	for i := range code {
		code[i] = imp.readUint32()
	}
	return code
}

func (imp *reader) readConstants() []interface{} {
	constants := make([]interface{}, imp.readUint32())
	for i := range constants {
		constants[i] = imp.readConstant()
	}
	return constants
}

func (imp *reader) readConstant() interface{} {
	switch imp.readByte() {
	case TagNil:
		return nil
	case TagBoolean:
		return imp.readByte() != 0
	case TagInteger:
		return imp.readLuaInteger()
	case TagNumber:
		return imp.readLuaNumber()
	case TagShortStr, TagLongStr:
		return imp.readString()
	default:
		panic("corrupted") // todo
	}
}

func (imp *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, imp.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: imp.readByte(),
			Idx:     imp.readByte(),
		}
	}
	return upvalues
}

func (imp *reader) readProtos(parentSource string) []*ProtoType {
	protos := make([]*ProtoType, imp.readUint32())
	for i := range protos {
		protos[i] = imp.readProto(parentSource)
	}
	return protos
}

func (imp *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, imp.readUint32())
	for i := range lineInfo {
		lineInfo[i] = imp.readUint32()
	}
	return lineInfo
}

func (imp *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, imp.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: imp.readString(),
			StartPC: imp.readUint32(),
			EndPC:   imp.readUint32(),
		}
	}
	return locVars
}

func (imp *reader) readUpvalueNames() []string {
	names := make([]string, imp.readUint32())
	for i := range names {
		names[i] = imp.readString()
	}
	return names
}
