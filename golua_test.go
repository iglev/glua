package glua

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/iglev/glua/api"
	"github.com/iglev/glua/compiler/lexer"
	"github.com/iglev/glua/state"
)

/*
func printHeader(f *binchunk.ProtoType) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *binchunk.ProtoType) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		// fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, c)
		i := vm.Instruction(c)
		fmt.Printf("\t%d\t[%s]\t%s \t", pc+1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

func printOperands(i vm.Instruction) {
	switch i.OpMode() {
	case vm.IABC:
		a, b, c := i.ABC()

		fmt.Printf("%d", a)
		if i.BMode() != vm.OpArgN {
			if b > 0xFF {
				fmt.Printf(" %d", -1-b&0xFF)
			} else {
				fmt.Printf(" %d", b)
			}
		}
		if i.CMode() != vm.OpArgN {
			if c > 0xFF {
				fmt.Printf(" %d", -1-c&0xFF)
			} else {
				fmt.Printf(" %d", c)
			}
		}
	case vm.IABx:
		a, bx := i.ABx()

		fmt.Printf("%d", a)
		if i.BMode() == vm.OpArgK {
			fmt.Printf(" %d", -1-bx)
		} else if i.BMode() == vm.OpArgU {
			fmt.Printf(" %d", bx)
		}
	case vm.IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case vm.IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
	}
}

func printDetail(f *binchunk.ProtoType) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, locVar.VarName, locVar.StartPC+1, locVar.EndPC+1)
	}

	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n",
			i, upvalName(f, i), upval.Instack, upval.Idx)
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(f *binchunk.ProtoType, idx int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[idx]
	}
	return "-"
}

func list(p *binchunk.ProtoType) {
	printHeader(p)
	printCode(p)
	printDetail(p)
	for _, subP := range p.Protos {
		list(subP)
	}
}

// TestBinChunk bin chunk
func TestBinChunk(t *testing.T) {
	data, err := ioutil.ReadFile("luac.out")
	if err != nil {
		panic(err)
	}
	proto := binchunk.Undump(data)
	list(proto)
}

func printStack(ls api.LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case api.LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case api.LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

// TestLuaState lua state
func TestLuaState(t *testing.T) {
	lstate := state.New(20, nil)

	lstate.PushBoolean(true)
	printStack(lstate)
	lstate.PushInteger(10)
	printStack(lstate)
	lstate.PushNil()
	printStack(lstate)
	lstate.PushString("hello")
	printStack(lstate)
	lstate.PushValue(-4)
	printStack(lstate)
	lstate.Replace(3)
	printStack(lstate)
	lstate.SetTop(6)
	printStack(lstate)
	lstate.Remove(-3)
	printStack(lstate)
	lstate.SetTop(-5)
	printStack(lstate)
}

// TestArith arith
func TestArith(t *testing.T) {
	lstate := state.New(20, nil)
	lstate.PushInteger(1)
	lstate.PushString("2.0")
	lstate.PushString("3.0")
	lstate.PushNumber(4.0)
	printStack(lstate)

	lstate.Arith(api.LUA_OPADD)
	printStack(lstate)
	lstate.Arith(api.LUA_OPBNOT)
	printStack(lstate)
	lstate.Len(2)
	printStack(lstate)
	lstate.Concat(3)
	printStack(lstate)
	lstate.PushBoolean(lstate.Compare(1, 2, api.LUA_OPEQ))
	printStack(lstate)
}

func luaMain(proto *binchunk.ProtoType) {
	nRegs := int(proto.MaxStackSize)
	lstate := state.New(nRegs+8, proto)
	lstate.SetTop(nRegs)
	for {
		pc := lstate.PC()
		inst := vm.Instruction(lstate.Fetch())
		if inst.Opcode() != vm.OP_RETURN {
			inst.Execute(lstate)

			fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
			printStack(lstate)
		} else {
			fmt.Printf("[%02d] %s \n", pc+1, inst.OpName())
			break
		}
	}
}

// TestVM vm
func TestVM(t *testing.T) {
	data, err := ioutil.ReadFile("luac.out")
	if err != nil {
		panic(err)
	}
	proto := binchunk.Undump(data)
	luaMain(proto)
}

// TestLoad load
func TestLoad(t *testing.T) {
	data, err := ioutil.ReadFile("luac.out")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Load(data, "luac.out", "b")
	ls.Call(0, 0)
}
*/

func print(ls api.LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}

func getMetatable(ls api.LuaState) int {
	if !ls.GetMetatable(1) {
		ls.PushNil()
	}
	return 1
}

func setMetatable(ls api.LuaState) int {
	ls.SetMetatable(1)
	return 1
}

func next(ls api.LuaState) int {
	ls.SetTop(2)
	if ls.Next(1) {
		return 2
	} else {
		ls.PushNil()
		return 1
	}
}

func pairs(ls api.LuaState) int {
	ls.PushGoFunction(next)
	ls.PushValue(1)
	ls.PushNil()
	return 3
}

func _iPairsAux(ls api.LuaState) int {
	i := ls.ToInteger(2) + 1
	ls.PushInteger(i)
	if ls.GetI(1, i) == api.LUA_TNIL {
		return 1
	}
	return 2
}

func ipairs(ls api.LuaState) int {
	ls.PushGoFunction(_iPairsAux)
	ls.PushValue(1)
	ls.PushNil()
	return 3
}

func errorFunc(ls api.LuaState) int {
	return ls.Error()
}

func pCall(ls api.LuaState) int {
	nArgs := ls.GetTop() - 1
	status := ls.PCall(nArgs, -1, 0)
	ls.PushBoolean(status == api.LUA_OK)
	ls.Insert(1)
	return ls.GetTop()
}

func typeFunc(ls api.LuaState) int {
	name := ls.TypeName(ls.Type(-1))
	ls.PushString(name)
	return 1
}

func testLexer(chunk, chunkName string) {
	lex := lexer.NewLexer(chunk, chunkName)
	for {
		line, kind, token := lex.NextToken()
		fmt.Printf("[%2d] [%-10s] %s\n",
			line, kindToCategory(kind), token)
		if kind == lexer.TOKEN_EOF {
			break
		}
	}
}

func kindToCategory(kind int) string {
	switch {
	case kind < lexer.TOKEN_SEP_SEMI:
		return "other"
	case kind <= lexer.TOKEN_SEP_RCURLY:
		return "separator"
	case kind <= lexer.TOKEN_OP_NOT:
		return "operator"
	case kind <= lexer.TOKEN_KW_WHILE:
		return "keyword"
	case kind == lexer.TOKEN_IDENTIFIER:
		return "identifier"
	case kind == lexer.TOKEN_NUMBER:
		return "number"
	case kind == lexer.TOKEN_STRING:
		return "string"
	default:
		return "other"
	}
}

// TestLexer lexer
func TestLexer(t *testing.T) {
	data, err := ioutil.ReadFile("test.lua")
	if err != nil {
		panic(err)
	}
	testLexer(string(data), "test.lua")
}

func testParser(chunk, chunkName string) {
	/*
		ast := parser.Parse(chunk, chunkName)
		b, err := json.Marshal(ast)
		if err != nil {
			panic(err)
		}
		println(string(b))
	*/
}

// TestParser parser
func TestParser(t *testing.T) {
	data, err := ioutil.ReadFile("test.lua")
	if err != nil {
		panic(err)
	}
	testParser(string(data), "test.lua")
}

// TestRegister register
func TestRegister(t *testing.T) {
	data, err := ioutil.ReadFile("test.lua")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Register("print", print)
	ls.Register("getmetatable", getMetatable)
	ls.Register("setmetatable", setMetatable)
	ls.Register("next", next)
	ls.Register("pairs", pairs)
	ls.Register("ipairs", ipairs)
	ls.Register("error", errorFunc)
	ls.Register("pcall", pCall)
	ls.Register("type", typeFunc)
	ls.Load(data, "test.lua", "b")
	ls.Call(0, 0)
}
