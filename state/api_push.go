package state

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
