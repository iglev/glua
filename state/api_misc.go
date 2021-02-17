package state

import "github.com/iglev/glua/number"

// Len - lua_len
func (l *luaState) Len(idx int) {
	val := l.stack.get(idx)

	if s, ok := val.(string); ok {
		l.stack.push(int64(len(s)))
	} else if result, ok := callMetamethod(val, val, "__len", l); ok {
		l.stack.push(result)
	} else if t, ok := val.(*luaTable); ok {
		l.stack.push(int64(t.len()))
	} else {
		panic("length error!")
	}
}

// Concat - lua_concat
func (l *luaState) Concat(n int) {
	if n == 0 {
		l.stack.push("")
	} else if n >= 2 {
		for i := 1; i < n; i++ {
			if l.IsString(-1) && l.IsString(-2) {
				s2 := l.ToString(-1)
				s1 := l.ToString(-2)
				l.stack.pop()
				l.stack.pop()
				l.stack.push(s1 + s2)
				continue
			}

			b := l.stack.pop()
			a := l.stack.pop()
			if result, ok := callMetamethod(a, b, "__concat", l); ok {
				l.stack.push(result)
				continue
			}

			panic("concatenation error!")
		}
	}
	// n == 1, do nothing
}

// Next - lua_next
func (l *luaState) Next(idx int) bool {
	val := l.stack.get(idx)
	if t, ok := val.(*luaTable); ok {
		key := l.stack.pop()
		if nextKey := t.nextKey(key); nextKey != nil {
			l.stack.push(nextKey)
			l.stack.push(t.get(nextKey))
			return true
		}
		return false
	}
	panic("table expected!")
}

// Error - lua_error
func (l *luaState) Error() int {
	err := l.stack.pop()
	panic(err)
}

// StringToNumber - lua_stringtonumber
func (l *luaState) StringToNumber(s string) bool {
	if n, ok := number.ParseInteger(s); ok {
		l.PushInteger(n)
		return true
	}
	if n, ok := number.ParseFloat(s); ok {
		l.PushNumber(n)
		return true
	}
	return false
}
