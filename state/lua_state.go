package state

import (
	"github.com/iglev/glua/binchunk"
)

type luaState struct {
	stack *luaStack
	proto *binchunk.ProtoType
	pc    int
}

// New new luaState
func New(stackSize int, proto *binchunk.ProtoType) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),
		proto: proto,
		pc:    0,
	}
}
