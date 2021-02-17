package binchunk

const (
	// LuaSignature lua binchunk header signature
	LuaSignature = "\x1bLua"
	// LuacVersion luac version (5.3 -> 5*16+3 -> 0x53)
	LuacVersion = 0x53
	// LuacFormat luac format
	LuacFormat = 0
	// LuacData luac data (0x19 93 0D 0A 1A 0A)
	LuacData = "\x19\x93\r\n\x1a\n"
	// CIntSize c int size
	CIntSize = 4
	// CSizetSize c size_t size
	CSizetSize = 8
	// InstructionSize instruction size
	InstructionSize = 4
	// LuaIntSize lua int size
	LuaIntSize = 8
	// LuaNumberSize lua number size
	LuaNumberSize = 8
	// LuacInt luac int (check big/little endian)
	LuacInt = 0x5678
	// LuacNum luac number (check double IEEE754)
	LuacNum = 370.5
)

const (
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

// BinaryChunk binary chunk
type BinaryChunk struct {
	Header
	sizeUpvalues byte
	mainFunc     *ProtoType
}

// Header binary chunk header
type Header struct {
	Signature       [4]byte
	Version         byte
	Format          byte
	LuacData        [6]byte
	CIntSize        byte
	SizetSize       byte
	InstructionSize byte
	LuaIntSize      byte
	LuaNumberSize   byte
	LuacInt         int64
	LuacNum         float64
}

// ProtoType prot type
type ProtoType struct {
	Source          string
	LineDefined     uint32
	LastLineDefined uint32
	NumParams       byte
	IsVararg        byte
	MaxStackSize    byte
	Code            []uint32
	Constants       []interface{}
	Upvalues        []Upvalue
	Protos          []*ProtoType
	LineInfo        []uint32
	LocVars         []LocVar
	UpvalueNames    []string
}

// Upvalue upvalue
type Upvalue struct {
	Instack byte
	Idx     byte
}

// LocVar local var
type LocVar struct {
	VarName string
	StartPC uint32
	EndPC   uint32
}

func IsBinaryChunk(data []byte) bool {
	return len(data) > 4 &&
		string(data[:4]) == LuaSignature
}

// Undump ...
func Undump(data []byte) *ProtoType {
	reader := &Reader{data}
	reader.CheckHeader()
	reader.ReadByte()
	return reader.ReadProto("")
}
