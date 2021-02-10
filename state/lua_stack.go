package state

type luaStack struct {
	slots []luaValue
	top   int
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

func (l *luaStack) check(n int) {
	free := len(l.slots) - l.top
	for i := free; i < n; i++ {
		l.slots = append(l.slots, nil)
	}
}

func (l *luaStack) push(val luaValue) {
	if l.top == len(l.slots) {
		panic("stack overflow!")
	}
	l.slots[l.top] = val
	l.top++
}

func (l *luaStack) pop() luaValue {
	if l.top < 1 {
		panic("stack underflow!")
	}
	l.top--
	val := l.slots[l.top]
	l.slots[l.top] = nil
	return val
}

func (l *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + l.top + 1
}

func (l *luaStack) isValid(idx int) bool {
	absIdx := l.absIndex(idx)
	return absIdx > 0 && absIdx <= l.top
}

func (l *luaStack) get(idx int) luaValue {
	absIdx := l.absIndex(idx)
	if absIdx > 0 && absIdx <= l.top {
		return l.slots[absIdx-1]
	}
	return nil
}

func (l *luaStack) set(idx int, val luaValue) {
	absIdx := l.absIndex(idx)
	if absIdx > 0 && absIdx <= l.top {
		l.slots[absIdx-1] = val
		return
	}
	panic("invalid index!")
}

func (l *luaStack) reverse(from, to int) {
	slots := l.slots
	for from < to {
		slots[from], slots[to] = slots[to], slots[from]
		from++
		to--
	}
}
