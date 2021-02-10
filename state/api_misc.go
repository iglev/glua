package state

// Len - lua_len
func (l *luaState) Len(idx int) {
	val := l.stack.get(idx)

	if s, ok := val.(string); ok {
		l.stack.push(int64(len(s)))
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

			panic("concatenation error!")
		}
	}
	// n == 1, do nothing
}
