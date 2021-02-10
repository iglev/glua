package state

// PC return pc
func (l *luaState) PC() int {
	return l.pc
}

// AddPC add pc
func (l *luaState) AddPC(n int) {
	l.pc += n
}

// Fetch fetch
func (l *luaState) Fetch() uint32 {
	i := l.proto.Code[l.pc]
	l.pc++
	return i
}

// GetConst get const
func (l *luaState) GetConst(idx int) {
	c := l.proto.Constants[idx]
	l.stack.push(c)
}

// GetRK get rk
func (l *luaState) GetRK(rk int) {
	if rk > 0xFF { // constant
		l.GetConst(rk & 0xFF)
	} else { // register
		l.PushValue(rk + 1)
	}
}
