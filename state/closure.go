package state

import (
	"github.com/iglev/glua/api"
	"github.com/iglev/glua/binchunk"
)

type closure struct {
	proto  *binchunk.ProtoType
	goFunc api.GoFunction
}

func newLuaClosure(proto *binchunk.ProtoType) *closure {
	return &closure{proto: proto}
}

func newGoClosure(f api.GoFunction) *closure {
	return &closure{goFunc: f}
}
