package state

import (
	"github.com/iglev/glua/api"
)

// CreateTable - lua_createtable
func (l *luaState) CreateTable(nArr, nRec int) {
	t := newLuaTable(nArr, nRec)
	l.stack.push(t)
}

// NewTable - lua_newtable
func (l *luaState) NewTable() {
	l.CreateTable(0, 0)
}

func (l *luaState) getTable(t, k luaValue, raw bool) api.LuaType {
	if tbl, ok := t.(*luaTable); ok {
		v := tbl.get(k)
		if raw || v != nil || !tbl.hasMetafield("__index") {
			l.stack.push(v)
			return typeOf(v)
		}
	}

	if !raw {
		if mf := getMetafield(t, "__index", l); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				return l.getTable(x, k, false)
			case *closure:
				l.stack.push(mf)
				l.stack.push(t)
				l.stack.push(k)
				l.Call(2, 1)
				v := l.stack.get(-1)
				return typeOf(v)
			}
		}
	}
	panic("not a table") // todo
}

// GetTable - lua_gettable
func (l *luaState) GetTable(idx int) api.LuaType {
	t := l.stack.get(idx)
	k := l.stack.pop()
	return l.getTable(t, k, false)
}

// GetField - lua_getfield
func (l *luaState) GetField(idx int, k string) api.LuaType {
	t := l.stack.get(idx)
	return l.getTable(t, k, false)
}

// GetI - lua_geti
func (l *luaState) GetI(idx int, i int64) api.LuaType {
	t := l.stack.get(idx)
	return l.getTable(t, i, false)
}

// RawGet - lua_rawget
func (l *luaState) RawGet(idx int) api.LuaType {
	t := l.stack.get(idx)
	k := l.stack.pop()
	return l.getTable(t, k, true)
}

// RawGetI - lua_rawgeti
func (l *luaState) RawGetI(idx int, i int64) api.LuaType {
	t := l.stack.get(idx)
	return l.getTable(t, i, true)
}

// GetGlobal - lua_getglobal
func (l *luaState) GetGlobal(name string) api.LuaType {
	t := l.registry.get(api.LUA_RIDX_GLOBALS)
	return l.getTable(t, name, false)
}

// GetMetatable - lua_getmetatable
func (l *luaState) GetMetatable(idx int) bool {
	val := l.stack.get(idx)
	if mt := getMetatable(val, l); mt != nil {
		l.stack.push(mt)
		return true
	}
	return false
}
