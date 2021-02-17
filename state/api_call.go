package state

import (
	"github.com/iglev/glua/api"
	"github.com/iglev/glua/binchunk"
	"github.com/iglev/glua/compiler"
	"github.com/iglev/glua/vm"
)

// Load - lua_load
func (l *luaState) Load(chunk []byte, chunkName, mode string) int {
	var proto *binchunk.ProtoType
	if binchunk.IsBinaryChunk(chunk) {
		proto = binchunk.Undump(chunk)
	} else {
		proto = compiler.Compile(string(chunk), chunkName)
	}

	c := newLuaClosure(proto)
	l.stack.push(c)
	if len(proto.Upvalues) > 0 {
		env := l.registry.get(api.LUA_RIDX_GLOBALS)
		c.upvals[0] = &upvalue{&env}
	}
	return api.LUA_OK
}

// Call - lua_call
func (l *luaState) Call(nArgs, nResults int) {
	val := l.stack.get(-(nArgs + 1))
	c, ok := val.(*closure)
	if !ok {
		if mf := getMetafield(val, "__call", l); mf != nil {
			if c, ok = mf.(*closure); ok {
				l.stack.push(val)
				l.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}

	if ok {
		if c.proto != nil {
			l.callLuaClosure(nArgs, nResults, c)
		} else {
			l.callGoClosure(nArgs, nResults, c)
		}
	} else {
		panic("not function")
	}
}

func (l *luaState) callGoClosure(nArgs, nResults int, c *closure) {
	// create new lua stack
	newStack := newLuaStack(nArgs+api.LUA_MINSTACK, l)
	newStack.closure = c

	// pass args, pop func
	if nArgs > 0 {
		args := l.stack.popN(nArgs)
		newStack.pushN(args, nArgs)
	}
	l.stack.pop()

	// run closure
	l.pushLuaStack(newStack)
	r := c.goFunc(l)
	l.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(r)
		l.stack.check(len(results))
		l.stack.pushN(results, nResults)
	}
}

func (l *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	nRegs := int(c.proto.MaxStackSize)
	nParams := int(c.proto.NumParams)
	isVararg := c.proto.IsVararg == 1

	// create new lua stack
	newStack := newLuaStack(nRegs+api.LUA_MINSTACK, l)
	newStack.closure = c

	// pass args, pop func
	funcAndArgs := l.stack.popN(nArgs + 1)
	newStack.pushN(funcAndArgs[1:], nParams)
	newStack.top = nRegs
	if nArgs > nParams && isVararg {
		newStack.varargs = funcAndArgs[nParams+1:]
	}

	// run closure
	l.pushLuaStack(newStack)
	l.runLuaClosure()
	l.popLuaStack()

	// return results
	if nResults != 0 {
		results := newStack.popN(newStack.top - nRegs)
		l.stack.check(len(results))
		l.stack.pushN(results, nResults)
	}
}

func (l *luaState) runLuaClosure() {
	for {
		inst := vm.Instruction(l.Fetch())
		inst.Execute(l)
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
	}
}

// PCall - lua_pcall
func (l *luaState) PCall(nArgs, nResults, msgh int) (status int) {
	caller := l.stack
	status = api.LUA_ERRRUN

	// catch error
	defer func() {
		if err := recover(); err != nil {
			if msgh != 0 {
				panic(err)
			}
			for l.stack != caller {
				l.popLuaStack()
			}
			l.stack.push(err)
		}
	}()

	l.Call(nArgs, nResults)
	status = api.LUA_OK
	return
}
