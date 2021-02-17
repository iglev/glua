package parser

import (
	"github.com/iglev/glua/compiler/ast"
	"github.com/iglev/glua/compiler/lexer"
)

// prefixexp ::= var | functioncall | ‘(’ exp ‘)’
// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args

/*
prefixexp ::= Name
	| ‘(’ exp ‘)’
	| prefixexp ‘[’ exp ‘]’
	| prefixexp ‘.’ Name
	| prefixexp [‘:’ Name] args
*/
func parsePrefixExp(lex *lexer.Lexer) ast.Exp {
	var exp ast.Exp
	if lex.LookAhead() == lexer.TOKEN_IDENTIFIER {
		line, name := lex.NextIdentifier() // Name
		exp = &ast.NameExp{line, name}
	} else { // ‘(’ exp ‘)’
		exp = parseParensExp(lex)
	}
	return _finishPrefixExp(lex, exp)
}

func parseParensExp(lex *lexer.Lexer) ast.Exp {
	lex.NextTokenOfKind(lexer.TOKEN_SEP_LPAREN) // (
	exp := parseExp(lex)                        // exp
	lex.NextTokenOfKind(lexer.TOKEN_SEP_RPAREN) // )

	switch exp.(type) {
	case *ast.VarargExp, *ast.FuncCallExp, *ast.NameExp, *ast.TableAccessExp:
		return &ast.ParensExp{exp}
	}

	// no need to keep parens
	return exp
}

func _finishPrefixExp(lex *lexer.Lexer, exp ast.Exp) ast.Exp {
	for {
		switch lex.LookAhead() {
		case lexer.TOKEN_SEP_LBRACK: // prefixexp ‘[’ exp ‘]’
			lex.NextToken()                             // ‘[’
			keyExp := parseExp(lex)                     // exp
			lex.NextTokenOfKind(lexer.TOKEN_SEP_RBRACK) // ‘]’
			exp = &ast.TableAccessExp{lex.Line(), exp, keyExp}
		case lexer.TOKEN_SEP_DOT: // prefixexp ‘.’ Name
			lex.NextToken()                    // ‘.’
			line, name := lex.NextIdentifier() // Name
			keyExp := &ast.StringExp{line, name}
			exp = &ast.TableAccessExp{line, exp, keyExp}
		case lexer.TOKEN_SEP_COLON, // prefixexp ‘:’ Name args
			lexer.TOKEN_SEP_LPAREN, lexer.TOKEN_SEP_LCURLY, lexer.TOKEN_STRING: // prefixexp args
			exp = _finishFuncCallExp(lex, exp)
		default:
			return exp
		}
	}
	return exp
}

// functioncall ::=  prefixexp args | prefixexp ‘:’ Name args
func _finishFuncCallExp(lex *lexer.Lexer, prefixExp ast.Exp) *ast.FuncCallExp {
	nameExp := _parseNameExp(lex)
	line := lex.Line() // todo
	args := _parseArgs(lex)
	lastLine := lex.Line()
	return &ast.FuncCallExp{line, lastLine, prefixExp, nameExp, args}
}

func _parseNameExp(lex *lexer.Lexer) *ast.StringExp {
	if lex.LookAhead() == lexer.TOKEN_SEP_COLON {
		lex.NextToken()
		line, name := lex.NextIdentifier()
		return &ast.StringExp{line, name}
	}
	return nil
}

// args ::=  ‘(’ [explist] ‘)’ | tableconstructor | LiteralString
func _parseArgs(lex *lexer.Lexer) (args []ast.Exp) {
	switch lex.LookAhead() {
	case lexer.TOKEN_SEP_LPAREN: // ‘(’ [explist] ‘)’
		lex.NextToken() // TOKEN_SEP_LPAREN
		if lex.LookAhead() != lexer.TOKEN_SEP_RPAREN {
			args = parseExpList(lex)
		}
		lex.NextTokenOfKind(lexer.TOKEN_SEP_RPAREN)
	case lexer.TOKEN_SEP_LCURLY: // ‘{’ [fieldlist] ‘}’
		args = []ast.Exp{parseTableConstructorExp(lex)}
	default: // LiteralString
		line, str := lex.NextTokenOfKind(lexer.TOKEN_STRING)
		args = []ast.Exp{&ast.StringExp{line, str}}
	}
	return
}
