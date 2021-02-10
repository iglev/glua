package state

/*
http://www.lua.org/manual/5.3/manual.html
*/

// GetTop - lua_gettop
func (l *luaState) GetTop() int {
	return l.stack.top
}

// AbsIndex - lua_absindex
func (l *luaState) AbsIndex(idx int) int {
	return l.stack.absIndex(idx)
}

// CheckStack - lua_checkstack
func (l *luaState) CheckStack(n int) bool {
	l.stack.check(n)
	return true // never fails
}

// Pop - lua_pop
func (l *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		l.stack.pop()
	}
}

// Copy - lua_copy
func (l *luaState) Copy(fromIdx, toIdx int) {
	val := l.stack.get(fromIdx)
	l.stack.set(toIdx, val)
}

// PushValue - PushValue
func (l *luaState) PushValue(idx int) {
	val := l.stack.get(idx)
	l.stack.push(val)
}

// Replace - lua_replace
func (l *luaState) Replace(idx int) {
	val := l.stack.pop()
	l.stack.set(idx, val)
}

// Insert - lua_insert
func (l *luaState) Insert(idx int) {
	l.Rotate(idx, 1)
}

// Remove - lua_remove
func (l *luaState) Remove(idx int) {
	l.Rotate(idx, -1)
	l.Pop(1)
}

/*
1 2 3 4
2 1 4 3
3 4 1 2
*/
// Rotate - lua_rotate
func (l *luaState) Rotate(idx, n int) {
	t := l.stack.top - 1
	p := l.stack.absIndex(idx) - 1
	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}
	l.stack.reverse(p, m)
	l.stack.reverse(m+1, t)
	l.stack.reverse(p, t)
}

func (l *luaState) SetTop(idx int) {
	newTop := l.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underflow!")
	}
	n := l.stack.top - newTop
	if n > 0 {
		for i := 0; i < n; i++ {
			l.stack.pop()
		}
	} else if n < 0 {
		for i := 0; i > n; i-- {
			l.stack.push(nil)
		}
	}
}
