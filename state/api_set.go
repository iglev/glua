package state

import "github.com/iglev/glua/api"

// SetTable - lua_settable
func (l *luaState) SetTable(idx int) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	k := l.stack.pop()
	l.setTable(t, k, v, false)
}

// SetField - lua_setfield
func (l *luaState) SetField(idx int, k string) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	l.setTable(t, k, v, false)
}

// SetI - lua_seti
func (l *luaState) SetI(idx int, i int64) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	l.setTable(t, i, v, false)
}

// RawSet - lua_rawset
func (l *luaState) RawSet(idx int) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	k := l.stack.pop()
	l.setTable(t, k, v, true)
}

// RawSetI - lua_rawseti
func (l *luaState) RawSetI(idx int, i int64) {
	t := l.stack.get(idx)
	v := l.stack.pop()
	l.setTable(t, i, v, true)
}

// t[k]=v
func (l *luaState) setTable(t, k, v luaValue, raw bool) {
	if tbl, ok := t.(*luaTable); ok {
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex") {
			tbl.put(k, v)
			return
		}
	}

	if !raw {
		if mf := getMetafield(t, "__newindex", l); mf != nil {
			switch x := mf.(type) {
			case *luaTable:
				l.setTable(x, k, v, false)
				return
			case *closure:
				l.stack.push(mf)
				l.stack.push(t)
				l.stack.push(k)
				l.stack.push(v)
				l.Call(3, 0)
				return
			}
		}
	}

	panic("index error!")
}

// SetGlobal - lua_setglobal
func (l *luaState) SetGlobal(name string) {
	t := l.registry.get(api.LUA_RIDX_GLOBALS)
	v := l.stack.pop()
	l.setTable(t, name, v, false)
}

// Register - lua_register
func (l *luaState) Register(name string, f api.GoFunction) {
	l.PushGoFunction(f)
	l.SetGlobal(name)
}

// SetMetatable - lua_setmetatable
func (l *luaState) SetMetatable(idx int) {
	val := l.stack.get(idx)
	mtVal := l.stack.pop()

	if mtVal == nil {
		setMetatable(val, nil, l)
	} else if mt, ok := mtVal.(*luaTable); ok {
		setMetatable(val, mt, l)
	} else {
		panic("table expected!") // todo
	}
}
