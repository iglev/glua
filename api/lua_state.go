package api

// LuaType lua type
type LuaType = int

// ArithOp arith op
type ArithOp = int

// CompareOp compare op
type CompareOp = int

// GoFunction go function
type GoFunction func(LuaState) int

// LuaState lua state interface
type LuaState interface {
	/* basic stack manipulation */
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromIdx, toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx, n int)
	SetTop(idx int)

	/* access functions (stack -> Go) */
	TypeName(tp LuaType) string
	Type(idx int) LuaType
	IsNone(idx int) bool
	IsNil(idx int) bool
	IsNoneOrNil(idx int) bool
	IsBoolean(idx int) bool
	IsInteger(idx int) bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	IsTable(idx int) bool
	IsThread(idx int) bool
	IsFunction(idx int) bool
	IsGoFunction(idx int) bool
	ToBoolean(idx int) bool
	ToInteger(idx int) int64
	ToIntegerX(idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64, bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)
	ToGoFunction(idx int) GoFunction

	/* push functions (Go -> stack) */
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	PushGoFunction(f GoFunction)
	PushGlobalTable()

	/* Comparison and arithmetic function */
	Arith(op ArithOp)
	Compare(idx1, idx2 int, op CompareOp) bool

	/* get functions (Lua -> stack) */
	NewTable()
	CreateTable(nArr, nRec int)
	GetTable(idx int) LuaType
	GetField(idx int, k string) LuaType
	GetI(idx int, i int64) LuaType
	GetGlobal(name string) LuaType

	/* set functions (stack -> Lua) */
	SetTable(idx int)
	SetField(idx int, k string)
	SetI(idx int, i int64)
	SetGlobal(name string)
	Register(name string, f GoFunction)

	/* 'load' and 'call' functions (load and run Lua code) */
	Load(chunk []byte, chunkName, mode string) int
	Call(nArgs, nResults int)

	/* miscellaneous functions */
	Len(idx int)
	Concat(n int)
}