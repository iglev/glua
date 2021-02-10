package state

import (
	"fmt"

	"github.com/iglev/glua/api"
)

// TypeName - lua_typename
func (l *luaState) TypeName(tp api.LuaType) string {
	switch tp {
	case api.LUA_TNONE:
		return "no value"
	case api.LUA_TNIL:
		return "nil"
	case api.LUA_TBOOLEAN:
		return "boolean"
	case api.LUA_TNUMBER:
		return "number"
	case api.LUA_TSTRING:
		return "string"
	case api.LUA_TTABLE:
		return "table"
	case api.LUA_TFUNCTION:
		return "function"
	case api.LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

// Type - lua_type
func (l *luaState) Type(idx int) api.LuaType {
	if l.stack.isValid(idx) {
		val := l.stack.get(idx)
		return typeOf(val)
	}
	return api.LUA_TNONE
}

// IsNone - lua_isnone
func (l *luaState) IsNone(idx int) bool {
	return l.Type(idx) == api.LUA_TNONE
}

// IsNil - lua_isnil
func (l *luaState) IsNil(idx int) bool {
	return l.Type(idx) == api.LUA_TNIL
}

// IsNoneOrNil - lua_isnoneornil
func (l *luaState) IsNoneOrNil(idx int) bool {
	return l.Type(idx) <= api.LUA_TNIL
}

// IsBoolean - lua_isboolean
func (l *luaState) IsBoolean(idx int) bool {
	return l.Type(idx) == api.LUA_TBOOLEAN
}

// IsTable - lua_istable
func (l *luaState) IsTable(idx int) bool {
	return l.Type(idx) == api.LUA_TTABLE
}

// IsFunction - lua_isfunction
func (l *luaState) IsFunction(idx int) bool {
	return l.Type(idx) == api.LUA_TFUNCTION
}

// IsThread - lua_isthread
func (l *luaState) IsThread(idx int) bool {
	return l.Type(idx) == api.LUA_TTHREAD
}

// IsString - lua_isstring
func (l *luaState) IsString(idx int) bool {
	t := l.Type(idx)
	return t == api.LUA_TSTRING || t == api.LUA_TNUMBER
}

// IsNumber - lua_isnumber
func (l *luaState) IsNumber(idx int) bool {
	_, ok := l.ToNumberX(idx)
	return ok
}

// IsInterger - lua_isinteger
func (l *luaState) IsInteger(idx int) bool {
	val := l.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

// ToBoolean - lua_toboolean
func (l *luaState) ToBoolean(idx int) bool {
	val := l.stack.get(idx)
	return convertToBoolean(val)
}

// ToInteger - lua_tointeger
func (l *luaState) ToInteger(idx int) int64 {
	i, _ := l.ToIntegerX(idx)
	return i
}

// ToIntegerX - lua_tointegerx
func (l *luaState) ToIntegerX(idx int) (int64, bool) {
	val := l.stack.get(idx)
	i, ok := val.(int64)
	return i, ok
}

// ToNumber - lua_tonumber
func (l *luaState) ToNumber(idx int) float64 {
	n, _ := l.ToNumberX(idx)
	return n
}

// ToNumberX - lua_tonumberx
func (l *luaState) ToNumberX(idx int) (float64, bool) {
	val := l.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

// ToString - lua_tostring
func (l *luaState) ToString(idx int) string {
	s, _ := l.ToStringX(idx)
	return s
}

// ToStringX - lua_tostringx
func (l *luaState) ToStringX(idx int) (string, bool) {
	val := l.stack.get(idx)

	switch x := val.(type) {
	case string:
		return x, true
	case int64, float64:
		s := fmt.Sprintf("%v", x) // todo
		l.stack.set(idx, s)
		return s, true
	default:
		return "", false
	}
}
