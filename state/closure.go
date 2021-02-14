package state

import (
	"github.com/iglev/glua/api"
	"github.com/iglev/glua/binchunk"
)

type upvalue struct {
	val *luaValue
}

type closure struct {
	proto  *binchunk.ProtoType
	goFunc api.GoFunction
	upvals []*upvalue
}

func newLuaClosure(proto *binchunk.ProtoType) *closure {
	c := &closure{proto: proto}
	if nUpvals := len(proto.Upvalues); nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}

func newGoClosure(f api.GoFunction, nUpvals int) *closure {
	c := &closure{goFunc: f}
	if nUpvals > 0 {
		c.upvals = make([]*upvalue, nUpvals)
	}
	return c
}
