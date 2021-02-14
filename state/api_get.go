package state

import "github.com/iglev/glua/api"

// CreateTable - lua_createtable
func (l *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	l.stack.push(t)
}

// NewTable - lua_newtable
func (l *luaState) NewTable() {
	l.CreateTable(0, 0)
}

func (l *luaState) getTable(t, k luaValue) api.LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		l.stack.push(v)
		return typeOf(v)
	}
	panic("not a table") // todo
}

// GetTable - lua_gettable
func (l *luaState) GetTable(idx int) api.LuaType {
	t := l.stack.get(idx)
	k := l.stack.pop()
	return l.getTable(t, k)
}

// GetField - lua_getfield
func (l *luaState) GetField(idx int, k string) api.LuaType {
	t := l.stack.get(idx)
	return l.getTable(t, k)
}

// GetI - lua_geti
func (l *luaState) GetI(idx int, i int64) api.LuaType {
	t := l.stack.get(idx)
	return l.getTable(t, i)
}

// GetGlobal - lua_getglobal
func (l *luaState) GetGlobal(name string) api.LuaType {
	t := l.registry.get(api.LUA_RIDX_GLOBALS)
	return l.getTable(t, name)
}
