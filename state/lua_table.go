package state

import (
	"math"

	"github.com/iglev/glua/number"
)

type luaTable struct {
	metatable *luaTable
	arr       []luaValue
	mp        map[luaValue]luaValue
	keys      map[luaValue]luaValue // used by next()
	lastKey   luaValue              // used by next()
	changed   bool                  // used by next()
}

func newLuaTable(nArr, nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0 {
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t.mp = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (lt *luaTable) hasMetafield(fieldName string) bool {
	return lt.metatable != nil &&
		lt.metatable.get(fieldName) != nil
}

func (lt *luaTable) len() int {
	return len(lt.arr)
}

func _floatToInteger(key luaValue) luaValue {
	if f, fok := key.(float64); fok {
		if i, ok := number.FloatToInteger(f); ok {
			return i
		}
	}
	return key
}

func (lt *luaTable) _shrinkArray() {
	for i := len(lt.arr) - 1; i >= 0; i-- {
		if lt.arr[i] == nil {
			lt.arr = lt.arr[0:i]
		}
	}
}

func (lt *luaTable) _expandArray() {
	for idx := int64(len(lt.arr)) + 1; true; idx++ {
		if val, found := lt.mp[idx]; found {
			delete(lt.mp, idx)
			lt.arr = append(lt.arr, val)
		} else {
			break
		}
	}
}

func (lt *luaTable) get(key luaValue) luaValue {
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok {
		if idx >= 1 && idx <= int64(len(lt.arr)) {
			return lt.arr[idx-1]
		}
	}
	return lt.mp[key]
}

func (lt *luaTable) put(key, val luaValue) {
	if key == nil {
		panic("table index is nil!")
	}
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is NaN")
	}

	lt.changed = true
	key = _floatToInteger(key)
	if idx, ok := key.(int64); ok && idx >= 1 {
		arrLen := int64(len(lt.arr))
		if idx <= arrLen {
			lt.arr[idx-1] = val
			if idx == arrLen && val == nil {
				lt._shrinkArray()
			}
			return
		}
		if idx == arrLen+1 {
			delete(lt.mp, key)
			if val != nil {
				lt.arr = append(lt.arr, val)
				lt._expandArray()
			}
			return
		}
	}
	if val != nil {
		if lt.mp == nil {
			lt.mp = make(map[luaValue]luaValue, 8)
		}
		lt.mp[key] = val
	} else {
		delete(lt.mp, key)
	}
}

func (lt *luaTable) initKeys() {
	lt.keys = make(map[luaValue]luaValue)
	var key luaValue = nil
	for i, v := range lt.arr {
		if v != nil {
			lt.keys[key] = int64(i + 1)
			key = int64(i + 1)
		}
	}
	for k, v := range lt.mp {
		if v != nil {
			lt.keys[key] = k
			key = k
		}
	}
	lt.lastKey = key
}

func (lt *luaTable) nextKey(key luaValue) luaValue {
	if lt.keys == nil || (key == nil && lt.changed) {
		lt.initKeys()
		lt.changed = false
	}

	nextKey := lt.keys[key]
	if nextKey == nil && key != nil && key != lt.lastKey {
		panic("invalid key to 'next'")
	}

	return nextKey
}
