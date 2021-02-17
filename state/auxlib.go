package state

import (
	"fmt"
	"io/ioutil"

	"github.com/iglev/glua/api"
	"github.com/iglev/glua/stdlib"
)

// Error2 - luaL_error
func (self *luaState) Error2(fmt string, a ...interface{}) int {
	self.PushFString(fmt, a...) // todo
	return self.Error()
}

// ArgError - luaL_argerror
func (self *luaState) ArgError(arg int, extraMsg string) int {
	// bad argument #arg to 'funcname' (extramsg)
	return self.Error2("bad argument #%d (%s)", arg, extraMsg) // todo
}

// CheckStack2 - luaL_checkstack
func (self *luaState) CheckStack2(sz int, msg string) {
	if !self.CheckStack(sz) {
		if msg != "" {
			self.Error2("stack overflow (%s)", msg)
		} else {
			self.Error2("stack overflow")
		}
	}
}

// ArgCheck - luaL_argcheck
func (self *luaState) ArgCheck(cond bool, arg int, extraMsg string) {
	if !cond {
		self.ArgError(arg, extraMsg)
	}
}

// CheckAny - luaL_checkany
func (self *luaState) CheckAny(arg int) {
	if self.Type(arg) == api.LUA_TNONE {
		self.ArgError(arg, "value expected")
	}
}

// CheckType - luaL_checktype
func (self *luaState) CheckType(arg int, t api.LuaType) {
	if self.Type(arg) != t {
		self.tagError(arg, t)
	}
}

// CheckInteger - luaL_checkinteger
func (self *luaState) CheckInteger(arg int) int64 {
	i, ok := self.ToIntegerX(arg)
	if !ok {
		self.intError(arg)
	}
	return i
}

// CheckNumber - luaL_checknumber
func (self *luaState) CheckNumber(arg int) float64 {
	f, ok := self.ToNumberX(arg)
	if !ok {
		self.tagError(arg, api.LUA_TNUMBER)
	}
	return f
}

// CheckString - luaL_checkstring & luaL_checklstring
func (self *luaState) CheckString(arg int) string {
	s, ok := self.ToStringX(arg)
	if !ok {
		self.tagError(arg, api.LUA_TSTRING)
	}
	return s
}

// OptInteger - luaL_optinteger
func (self *luaState) OptInteger(arg int, def int64) int64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckInteger(arg)
}

// OptNumber - luaL_optnumber
func (self *luaState) OptNumber(arg int, def float64) float64 {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckNumber(arg)
}

// OptString - luaL_optstring
func (self *luaState) OptString(arg int, def string) string {
	if self.IsNoneOrNil(arg) {
		return def
	}
	return self.CheckString(arg)
}

// DoFile - luaL_dofile
func (self *luaState) DoFile(filename string) bool {
	return self.LoadFile(filename) != api.LUA_OK ||
		self.PCall(0, api.LUA_MULTRET, 0) != api.LUA_OK
}

// DoString - luaL_dostring
func (self *luaState) DoString(str string) bool {
	return self.LoadString(str) != api.LUA_OK ||
		self.PCall(0, api.LUA_MULTRET, 0) != api.LUA_OK
}

// LoadFile - luaL_loadfile
func (self *luaState) LoadFile(filename string) int {
	return self.LoadFileX(filename, "bt")
}

// LoadFileX - luaL_loadfilex
func (self *luaState) LoadFileX(filename, mode string) int {
	if data, err := ioutil.ReadFile(filename); err == nil {
		return self.Load(data, "@"+filename, mode)
	}
	return api.LUA_ERRFILE
}

// LoadString - luaL_loadstring
func (self *luaState) LoadString(s string) int {
	return self.Load([]byte(s), s, "bt")
}

// TypeName2 - luaL_typename
func (self *luaState) TypeName2(idx int) string {
	return self.TypeName(self.Type(idx))
}

// Len2 - luaL_len
func (self *luaState) Len2(idx int) int64 {
	self.Len(idx)
	i, isNum := self.ToIntegerX(-1)
	if !isNum {
		self.Error2("object length is not an integer")
	}
	self.Pop(1)
	return i
}

// ToString2 - luaL_tolstring
func (self *luaState) ToString2(idx int) string {
	if self.CallMeta(idx, "__tostring") { /* metafield? */
		if !self.IsString(-1) {
			self.Error2("'__tostring' must return a string")
		}
	} else {
		switch self.Type(idx) {
		case api.LUA_TNUMBER:
			if self.IsInteger(idx) {
				self.PushString(fmt.Sprintf("%d", self.ToInteger(idx))) // todo
			} else {
				self.PushString(fmt.Sprintf("%g", self.ToNumber(idx))) // todo
			}
		case api.LUA_TSTRING:
			self.PushValue(idx)
		case api.LUA_TBOOLEAN:
			if self.ToBoolean(idx) {
				self.PushString("true")
			} else {
				self.PushString("false")
			}
		case api.LUA_TNIL:
			self.PushString("nil")
		default:
			tt := self.GetMetafield(idx, "__name") /* try name */
			var kind string
			if tt == api.LUA_TSTRING {
				kind = self.CheckString(-1)
			} else {
				kind = self.TypeName2(idx)
			}

			self.PushString(fmt.Sprintf("%s: %p", kind, self.ToPointer(idx)))
			if tt != api.LUA_TNIL {
				self.Remove(-2) /* remove '__name' */
			}
		}
	}
	return self.CheckString(-1)
}

// GetSubTable - luaL_getsubtable
func (self *luaState) GetSubTable(idx int, fname string) bool {
	if self.GetField(idx, fname) == api.LUA_TTABLE {
		return true /* table already there */
	}
	self.Pop(1) /* remove previous result */
	idx = self.stack.absIndex(idx)
	self.NewTable()
	self.PushValue(-1)        /* copy to be left at top */
	self.SetField(idx, fname) /* assign new table to field */
	return false              /* false, because did not find table there */
}

// GetMetafield - luaL_getmetafield
func (self *luaState) GetMetafield(obj int, event string) api.LuaType {
	if !self.GetMetatable(obj) { /* no metatable? */
		return api.LUA_TNIL
	}

	self.PushString(event)
	tt := self.RawGet(-2)
	if tt == api.LUA_TNIL { /* is metafield nil? */
		self.Pop(2) /* remove metatable and metafield */
	} else {
		self.Remove(-2) /* remove only metatable */
	}
	return tt /* return metafield type */
}

// CallMeta - luaL_callmeta
func (self *luaState) CallMeta(obj int, event string) bool {
	obj = self.AbsIndex(obj)
	if self.GetMetafield(obj, event) == api.LUA_TNIL { /* no metafield? */
		return false
	}

	self.PushValue(obj)
	self.Call(1, 1)
	return true
}

// OpenLibs - luaL_openlibs
func (self *luaState) OpenLibs() {
	libs := map[string]api.GoFunction{
		"_G":        stdlib.OpenBaseLib,
		"math":      stdlib.OpenMathLib,
		"table":     stdlib.OpenTableLib,
		"string":    stdlib.OpenStringLib,
		"utf8":      stdlib.OpenUTF8Lib,
		"os":        stdlib.OpenOSLib,
		"package":   stdlib.OpenPackageLib,
		"coroutine": stdlib.OpenCoroutineLib,
	}

	for name, fun := range libs {
		self.RequireF(name, fun, true)
		self.Pop(1)
	}
}

// RequireF - luaL_requiref
func (self *luaState) RequireF(modname string, openf api.GoFunction, glb bool) {
	self.GetSubTable(api.LUA_REGISTRYINDEX, "_LOADED")
	self.GetField(-1, modname) /* LOADED[modname] */
	if !self.ToBoolean(-1) {   /* package not already loaded? */
		self.Pop(1) /* remove field */
		self.PushGoFunction(openf)
		self.PushString(modname)   /* argument to open function */
		self.Call(1, 1)            /* call 'openf' to open module */
		self.PushValue(-1)         /* make copy of module (call result) */
		self.SetField(-3, modname) /* _LOADED[modname] = module */
	}
	self.Remove(-2) /* remove _LOADED table */
	if glb {
		self.PushValue(-1)      /* copy of module */
		self.SetGlobal(modname) /* _G[modname] = module */
	}
}

// NewLib - luaL_newlib
func (self *luaState) NewLib(l api.FuncReg) {
	self.NewLibTable(l)
	self.SetFuncs(l, 0)
}

// NewLibTable - luaL_newlibtable
func (self *luaState) NewLibTable(l api.FuncReg) {
	self.CreateTable(0, len(l))
}

// SetFuncs - luaL_setfuncs
func (self *luaState) SetFuncs(l api.FuncReg, nup int) {
	self.CheckStack2(nup, "too many upvalues")
	for name, fun := range l { /* fill the table with given functions */
		for i := 0; i < nup; i++ { /* copy upvalues to the top */
			self.PushValue(-nup)
		}
		// r[-(nup+2)][name]=fun
		self.PushGoClosure(fun, nup) /* closure with those upvalues */
		self.SetField(-(nup + 2), name)
	}
	self.Pop(nup) /* remove upvalues */
}

func (self *luaState) intError(arg int) {
	if self.IsNumber(arg) {
		self.ArgError(arg, "number has no integer representation")
	} else {
		self.tagError(arg, api.LUA_TNUMBER)
	}
}

func (self *luaState) tagError(arg int, tag api.LuaType) {
	self.typeError(arg, self.TypeName(api.LuaType(tag)))
}

func (self *luaState) typeError(arg int, tname string) int {
	var typeArg string /* name for the type of the actual argument */
	if self.GetMetafield(arg, "__name") == api.LUA_TSTRING {
		typeArg = self.ToString(-1) /* use the given type name */
	} else if self.Type(arg) == api.LUA_TLIGHTUSERDATA {
		typeArg = "light userdata" /* special name for messages */
	} else {
		typeArg = self.TypeName2(arg) /* standard name */
	}
	msg := tname + " expected, got " + typeArg
	self.PushString(msg)
	return self.ArgError(arg, msg)
}
