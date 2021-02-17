package compiler

import (
	"github.com/iglev/glua/binchunk"
	"github.com/iglev/glua/compiler/codegen"
	"github.com/iglev/glua/compiler/parser"
)

func Compile(chunk, chunkName string) *binchunk.ProtoType {
	ast := parser.Parse(chunk, chunkName)
	return codegen.GenProto(ast)
}
