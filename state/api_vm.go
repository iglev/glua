package state

// PC return pc
func (l *luaState) PC() int {
	return l.stack.pc
}

// AddPC add pc
func (l *luaState) AddPC(n int) {
	l.stack.pc += n
}

// Fetch fetch
func (l *luaState) Fetch() uint32 {
	i := l.stack.closure.proto.Code[l.stack.pc]
	l.stack.pc++
	return i
}

// GetConst get const
func (l *luaState) GetConst(idx int) {
	c := l.stack.closure.proto.Constants[idx]
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

// RegisterCount register count
func (l *luaState) RegisterCount() int {
	return int(l.stack.closure.proto.MaxStackSize)
}

// LoadVararg load vararg
func (l *luaState) LoadVararg(n int) {
	if n < 0 {
		n = len(l.stack.varargs)
	}

	l.stack.check(n)
	l.stack.pushN(l.stack.varargs, n)
}

// LoadProto load proto
func (l *luaState) LoadProto(idx int) {
	stack := l.stack
	subProto := stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)
	for i, uvInfo := range subProto.Upvalues {
		uvIdx := int(uvInfo.Idx)
		if uvInfo.Instack == 1 {
			if stack.openuvs == nil {
				stack.openuvs = make(map[int]*upvalue)
			}
			if openuv, found := stack.openuvs[uvIdx]; found {
				closure.upvals[i] = openuv
			} else {
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
		} else {
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}

func (l *luaState) CloseUpvalues(a int) {
	for i, openuv := range l.stack.openuvs {
		if i >= a-1 {
			val := *openuv.val
			openuv.val = &val
			delete(l.stack.openuvs, i)
		}
	}
}
