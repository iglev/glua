package codegen

import (
	"github.com/iglev/glua/binchunk"
	"github.com/iglev/glua/compiler/ast"
)

func GenProto(chunk *ast.Block) *binchunk.ProtoType {
	fd := &ast.FuncDefExp{
		IsVararg: true,
		Block:    chunk,
	}

	fi := newFuncInfo(nil, fd)
	fi.addLocVar("_ENV")
	cgFuncDefExp(fi, fd, 0)
	return toProto(fi.subFuncs[0])
}
