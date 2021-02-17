package state

import (
	"github.com/iglev/glua/api"
)

type luaState struct {
	registry *luaTable
	stack    *luaStack

	/* coroutine */
	coStatus int
	coCaller *luaState
	coChan   chan int
}

// New new luaState
func New() *luaState {
	ls := &luaState{}

	registry := newLuaTable(8, 0)
	registry.put(api.LUA_RIDX_MAINTHREAD, ls)
	registry.put(api.LUA_RIDX_GLOBALS, newLuaTable(0, 20))

	ls.registry = registry
	ls.pushLuaStack(newLuaStack(api.LUA_MINSTACK, ls))
	return ls
}

func (l *luaState) isMainThread() bool {
	return l.registry.get(api.LUA_RIDX_MAINTHREAD) == l
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
