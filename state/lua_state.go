package state

import "github.com/iglev/glua/api"

type luaState struct {
	registry *luaTable
	stack    *luaStack
}

// New new luaState
func New() *luaState {
	// new registry
	registry := newLuaTable(0, 0)
	registry.put(api.LUA_RIDX_GLOBALS, newLuaTable(0, 0))

	// new state
	ls := &luaState{registry: registry}
	ls.pushLuaStack(newLuaStack(api.LUA_MINSTACK, ls))

	return ls
}

func (l *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = l.stack
	l.stack = stack
}

func (l *luaState) popLuaStack() {
	stack := l.stack
	l.stack = stack.prev
	stack.prev = nil
}
