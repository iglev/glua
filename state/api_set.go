package state

import "github.com/iglev/glua/api"

// SetTable - lua_settable
func (l *luaState) SetTable(idx int) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	k := l.stack.pop()
	l.setTable(t, k, v)
}

// SetField - lua_setfield
func (l *luaState) SetField(idx int, k string) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	l.setTable(t, k, v)
}

// SetI - lua_seti
func (l *luaState) SetI(idx int, i int64) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	l.setTable(t, i, v)
}

// t[k]=v
func (l *luaState) setTable(t, k, v luaValue) {
	if tbl, ok := t.(*luaTable); ok {
		tbl.put(k, v)
		return
	}

	panic("not a table!")
}

// SetGlobal - lua_setglobal
func (l *luaState) SetGlobal(name string) {
	t := l.registry.get(api.LUA_RIDX_GLOBALS)
	v := l.stack.pop()
	l.setTable(t, name, v)
}

// Register - lua_register
func (l *luaState) Register(name string, f api.GoFunction) {
	l.PushGoFunction(f)
	l.SetGlobal(name)
}
