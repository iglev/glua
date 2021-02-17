package parser

import (
	"github.com/iglev/glua/compiler/ast"
	"github.com/iglev/glua/compiler/lexer"
	"github.com/iglev/glua/number"
)

// explist ::= exp {‘,’ exp}
func parseExpList(lex *lexer.Lexer) []ast.Exp {
	exps := make([]ast.Exp, 0, 4)
	exps = append(exps, parseExp(lex))
	for lex.LookAhead() == lexer.TOKEN_SEP_COMMA {
		lex.NextToken()
		exps = append(exps, parseExp(lex))
	}
	return exps
}

/*
exp ::=  nil | false | true | Numeral | LiteralString | ‘...’ | functiondef |
	 prefixexp | tableconstructor | exp binop exp | unop exp
*/
/*
exp   ::= exp12
exp12 ::= exp11 {or exp11}
exp11 ::= exp10 {and exp10}
exp10 ::= exp9 {(‘<’ | ‘>’ | ‘<=’ | ‘>=’ | ‘~=’ | ‘==’) exp9}
exp9  ::= exp8 {‘|’ exp8}
exp8  ::= exp7 {‘~’ exp7}
exp7  ::= exp6 {‘&’ exp6}
exp6  ::= exp5 {(‘<<’ | ‘>>’) exp5}
exp5  ::= exp4 {‘..’ exp4}
exp4  ::= exp3 {(‘+’ | ‘-’) exp3}
exp3  ::= exp2 {(‘*’ | ‘/’ | ‘//’ | ‘%’) exp2}
exp2  ::= {(‘not’ | ‘#’ | ‘-’ | ‘~’)} exp1
exp1  ::= exp0 {‘^’ exp2}
exp0  ::= nil | false | true | Numeral | LiteralString
		| ‘...’ | functiondef | prefixexp | tableconstructor
*/
func parseExp(lex *lexer.Lexer) ast.Exp {
	return parseExp12(lex)
}

// x or y
func parseExp12(lex *lexer.Lexer) ast.Exp {
	exp := parseExp11(lex)
	for lex.LookAhead() == lexer.TOKEN_OP_OR {
		line, op, _ := lex.NextToken()
		lor := &ast.BinopExp{line, op, exp, parseExp11(lex)}
		exp = optimizeLogicalOr(lor)
	}
	return exp
}

// x and y
func parseExp11(lex *lexer.Lexer) ast.Exp {
	exp := parseExp10(lex)
	for lex.LookAhead() == lexer.TOKEN_OP_AND {
		line, op, _ := lex.NextToken()
		land := &ast.BinopExp{line, op, exp, parseExp10(lex)}
		exp = optimizeLogicalAnd(land)
	}
	return exp
}

// compare
func parseExp10(lex *lexer.Lexer) ast.Exp {
	exp := parseExp9(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TOKEN_OP_LT, lexer.TOKEN_OP_GT, lexer.TOKEN_OP_NE,
			lexer.TOKEN_OP_LE, lexer.TOKEN_OP_GE, lexer.TOKEN_OP_EQ:
			line, op, _ := lex.NextToken()
			exp = &ast.BinopExp{line, op, exp, parseExp9(lex)}
		default:
			return exp
		}
	}
	return exp
}

// x | y
func parseExp9(lex *lexer.Lexer) ast.Exp {
	exp := parseExp8(lex)
	for lex.LookAhead() == lexer.TOKEN_OP_BOR {
		line, op, _ := lex.NextToken()
		bor := &ast.BinopExp{line, op, exp, parseExp8(lex)}
		exp = optimizeBitwiseBinaryOp(bor)
	}
	return exp
}

// x ~ y
func parseExp8(lex *lexer.Lexer) ast.Exp {
	exp := parseExp7(lex)
	for lex.LookAhead() == lexer.TOKEN_OP_BXOR {
		line, op, _ := lex.NextToken()
		bxor := &ast.BinopExp{line, op, exp, parseExp7(lex)}
		exp = optimizeBitwiseBinaryOp(bxor)
	}
	return exp
}

// x & y
func parseExp7(lex *lexer.Lexer) ast.Exp {
	exp := parseExp6(lex)
	for lex.LookAhead() == lexer.TOKEN_OP_BAND {
		line, op, _ := lex.NextToken()
		band := &ast.BinopExp{line, op, exp, parseExp6(lex)}
		exp = optimizeBitwiseBinaryOp(band)
	}
	return exp
}

// shift
func parseExp6(lex *lexer.Lexer) ast.Exp {
	exp := parseExp5(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TOKEN_OP_SHL, lexer.TOKEN_OP_SHR:
			line, op, _ := lex.NextToken()
			shx := &ast.BinopExp{line, op, exp, parseExp5(lex)}
			exp = optimizeBitwiseBinaryOp(shx)
		default:
			return exp
		}
	}
	return exp
}

// a .. b
func parseExp5(lex *lexer.Lexer) ast.Exp {
	exp := parseExp4(lex)
	if lex.LookAhead() != lexer.TOKEN_OP_CONCAT {
		return exp
	}

	line := 0
	exps := []ast.Exp{exp}
	for lex.LookAhead() == lexer.TOKEN_OP_CONCAT {
		line, _, _ = lex.NextToken()
		exps = append(exps, parseExp4(lex))
	}
	return &ast.ConcatExp{line, exps}
}

// x +/- y
func parseExp4(lex *lexer.Lexer) ast.Exp {
	exp := parseExp3(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TOKEN_OP_ADD, lexer.TOKEN_OP_SUB:
			line, op, _ := lex.NextToken()
			arith := &ast.BinopExp{line, op, exp, parseExp3(lex)}
			exp = optimizeArithBinaryOp(arith)
		default:
			return exp
		}
	}
	return exp
}

// *, %, /, //
func parseExp3(lex *lexer.Lexer) ast.Exp {
	exp := parseExp2(lex)
	for {
		switch lex.LookAhead() {
		case lexer.TOKEN_OP_MUL, lexer.TOKEN_OP_MOD, lexer.TOKEN_OP_DIV, lexer.TOKEN_OP_IDIV:
			line, op, _ := lex.NextToken()
			arith := &ast.BinopExp{line, op, exp, parseExp2(lex)}
			exp = optimizeArithBinaryOp(arith)
		default:
			return exp
		}
	}
	return exp
}

// unary
func parseExp2(lex *lexer.Lexer) ast.Exp {
	switch lex.LookAhead() {
	case lexer.TOKEN_OP_UNM, lexer.TOKEN_OP_BNOT, lexer.TOKEN_OP_LEN, lexer.TOKEN_OP_NOT:
		line, op, _ := lex.NextToken()
		exp := &ast.UnopExp{line, op, parseExp2(lex)}
		return optimizeUnaryOp(exp)
	}
	return parseExp1(lex)
}

// x ^ y
func parseExp1(lex *lexer.Lexer) ast.Exp { // pow is right associative
	exp := parseExp0(lex)
	if lex.LookAhead() == lexer.TOKEN_OP_POW {
		line, op, _ := lex.NextToken()
		exp = &ast.BinopExp{line, op, exp, parseExp2(lex)}
	}
	return optimizePow(exp)
}

func parseExp0(lex *lexer.Lexer) ast.Exp {
	switch lex.LookAhead() {
	case lexer.TOKEN_VARARG: // ...
		line, _, _ := lex.NextToken()
		return &ast.VarargExp{line}
	case lexer.TOKEN_KW_NIL: // nil
		line, _, _ := lex.NextToken()
		return &ast.NilExp{line}
	case lexer.TOKEN_KW_TRUE: // true
		line, _, _ := lex.NextToken()
		return &ast.TrueExp{line}
	case lexer.TOKEN_KW_FALSE: // false
		line, _, _ := lex.NextToken()
		return &ast.FalseExp{line}
	case lexer.TOKEN_STRING: // LiteralString
		line, _, token := lex.NextToken()
		return &ast.StringExp{line, token}
	case lexer.TOKEN_NUMBER: // Numeral
		return parseNumberExp(lex)
	case lexer.TOKEN_SEP_LCURLY: // tableconstructor
		return parseTableConstructorExp(lex)
	case lexer.TOKEN_KW_FUNCTION: // functiondef
		lex.NextToken()
		return parseFuncDefExp(lex)
	default: // prefixexp
		return parsePrefixExp(lex)
	}
}

func parseNumberExp(lex *lexer.Lexer) ast.Exp {
	line, _, token := lex.NextToken()
	if i, ok := number.ParseInteger(token); ok {
		return &ast.IntegerExp{line, i}
	} else if f, ok := number.ParseFloat(token); ok {
		return &ast.FloatExp{line, f}
	} else { // todo
		panic("not a number: " + token)
	}
}

// functiondef ::= function funcbody
// funcbody ::= ‘(’ [parlist] ‘)’ block end
func parseFuncDefExp(lex *lexer.Lexer) *ast.FuncDefExp {
	line := lex.Line()                                     // function
	lex.NextTokenOfKind(lexer.TOKEN_SEP_LPAREN)            // (
	parList, isVararg := _parseParList(lex)                // [parlist]
	lex.NextTokenOfKind(lexer.TOKEN_SEP_RPAREN)            // )
	block := parseBlock(lex)                               // block
	lastLine, _ := lex.NextTokenOfKind(lexer.TOKEN_KW_END) // end
	return &ast.FuncDefExp{line, lastLine, parList, isVararg, block}
}

// [parlist]
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
func _parseParList(lex *lexer.Lexer) (names []string, isVararg bool) {
	switch lex.LookAhead() {
	case lexer.TOKEN_SEP_RPAREN:
		return nil, false
	case lexer.TOKEN_VARARG:
		lex.NextToken()
		return nil, true
	}

	_, name := lex.NextIdentifier()
	names = append(names, name)
	for lex.LookAhead() == lexer.TOKEN_SEP_COMMA {
		lex.NextToken()
		if lex.LookAhead() == lexer.TOKEN_IDENTIFIER {
			_, name := lex.NextIdentifier()
			names = append(names, name)
		} else {
			lex.NextTokenOfKind(lexer.TOKEN_VARARG)
			isVararg = true
			break
		}
	}
	return
}

// tableconstructor ::= ‘{’ [fieldlist] ‘}’
func parseTableConstructorExp(lex *lexer.Lexer) *ast.TableConstructorExp {
	line := lex.Line()
	lex.NextTokenOfKind(lexer.TOKEN_SEP_LCURLY) // {
	keyExps, valExps := _parseFieldList(lex)    // [fieldlist]
	lex.NextTokenOfKind(lexer.TOKEN_SEP_RCURLY) // }
	lastLine := lex.Line()
	return &ast.TableConstructorExp{line, lastLine, keyExps, valExps}
}

// fieldlist ::= field {fieldsep field} [fieldsep]
func _parseFieldList(lex *lexer.Lexer) (ks, vs []ast.Exp) {
	if lex.LookAhead() != lexer.TOKEN_SEP_RCURLY {
		k, v := _parseField(lex)
		ks = append(ks, k)
		vs = append(vs, v)

		for _isFieldSep(lex.LookAhead()) {
			lex.NextToken()
			if lex.LookAhead() != lexer.TOKEN_SEP_RCURLY {
				k, v := _parseField(lex)
				ks = append(ks, k)
				vs = append(vs, v)
			} else {
				break
			}
		}
	}
	return
}

// fieldsep ::= ‘,’ | ‘;’
func _isFieldSep(tokenKind int) bool {
	return tokenKind == lexer.TOKEN_SEP_COMMA || tokenKind == lexer.TOKEN_SEP_SEMI
}

// field ::= ‘[’ exp ‘]’ ‘=’ exp | Name ‘=’ exp | exp
func _parseField(lex *lexer.Lexer) (k, v ast.Exp) {
	if lex.LookAhead() == lexer.TOKEN_SEP_LBRACK {
		lex.NextToken()                             // [
		k = parseExp(lex)                           // exp
		lex.NextTokenOfKind(lexer.TOKEN_SEP_RBRACK) // ]
		lex.NextTokenOfKind(lexer.TOKEN_OP_ASSIGN)  // =
		v = parseExp(lex)                           // exp
		return
	}

	exp := parseExp(lex)
	if nameExp, ok := exp.(*ast.NameExp); ok {
		if lex.LookAhead() == lexer.TOKEN_OP_ASSIGN {
			// Name ‘=’ exp => ‘[’ LiteralString ‘]’ = exp
			lex.NextToken()
			k = &ast.StringExp{nameExp.Line, nameExp.Name}
			v = parseExp(lex)
			return
		}
	}

	return nil, exp
}
