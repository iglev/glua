package state

import (
	"github.com/iglev/glua/api"
)

// RawEqual - lua_rawequal
func (l *luaState) RawEqual(idx1, idx2 int) bool {
	if !l.stack.isValid(idx1) || !l.stack.isValid(idx2) {
		return false
	}

	a := l.stack.get(idx1)
	b := l.stack.get(idx2)
	return _eq(a, b, nil)
}

// Compare - lua_compare
func (l *luaState) Compare(idx1, idx2 int, op api.CompareOp) bool {
	if !l.stack.isValid(idx1) || !l.stack.isValid(idx2) {
		return false
	}

	a := l.stack.get(idx1)
	b := l.stack.get(idx2)
	switch op {
	case api.LUA_OPEQ:
		return _eq(a, b, l)
	case api.LUA_OPLT:
		return _lt(a, b, l)
	case api.LUA_OPLE:
		return _le(a, b, l)
	default:
		panic("invalid compare op!")
	}
}

func _eq(a, b luaValue, l *luaState) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int64:
		switch y := b.(type) {
		case int64:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int64:
			return x == float64(y)
		default:
			return false
		}
	case *luaTable:
		if y, ok := b.(*luaTable); ok && x != y && l != nil {
			if result, ok := callMetamethod(x, y, "__eq", l); ok {
				return convertToBoolean(result)
			}
		}
		return a == b
	default:
		return a == b
	}
}

func _lt(a, b luaValue, l *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x < y
		case float64:
			return float64(x) < y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y
		case int64:
			return x < float64(y)
		}
	}
	if result, ok := callMetamethod(a, b, "__lt", l); ok {
		return convertToBoolean(result)
	}
	panic("comparison error!")
}

func _le(a, b luaValue, l *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	if result, ok := callMetamethod(a, b, "__le", l); ok {
		return convertToBoolean(result)
	} else if result, ok := callMetamethod(b, a, "__lt", l); ok {
		return !convertToBoolean(result)
	}
	panic("comparison error!")
}
