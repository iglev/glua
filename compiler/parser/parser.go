package parser

import (
	"github.com/iglev/glua/compiler/ast"
	"github.com/iglev/glua/compiler/lexer"
)

/* recursive descent parser */

func Parse(chunk, chunkName string) *ast.Block {
	lex := lexer.NewLexer(chunk, chunkName)
	block := parseBlock(lex)
	lex.NextTokenOfKind(lexer.TOKEN_EOF)
	return block
}
