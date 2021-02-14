package state

import "github.com/iglev/glua/api"

/*
http://www.lua.org/manual/5.3/manual.html
*/

// PushNil - lua_pushnil
func (l *luaState) PushNil() {
	l.stack.push(nil)
}

// PushBoolean - lua_pushboolean
func (l *luaState) PushBoolean(b bool) {
	l.stack.push(b)
}

// PushInteger - lua_pushinteger
func (l *luaState) PushInteger(n int64) {
	l.stack.push(n)
}

// PushNumber - lua_pushnumber
func (l *luaState) PushNumber(n float64) {
	l.stack.push(n)
}

// PushString - lua_pushstring
func (l *luaState) PushString(s string) {
	l.stack.push(s)
}

// PushGoFunction - lua_pushcfunction
func (l *luaState) PushGoFunction(f api.GoFunction) {
	l.stack.push(newGoClosure(f))
}

// PushGlobalTable - lua_pushglobaltable
func (l *luaState) PushGlobalTable() {
	global := l.registry.get(api.LUA_RIDX_GLOBALS)
	l.stack.push(global)
}
