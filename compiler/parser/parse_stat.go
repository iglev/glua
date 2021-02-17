package parser

import (
	"github.com/iglev/glua/compiler/ast"
	"github.com/iglev/glua/compiler/lexer"
)

var _statEmpty = &ast.EmptyStat{}

/*
stat ::=  ‘;’
	| break
	| ‘::’ Name ‘::’
	| goto Name
	| do block end
	| while exp do block end
	| repeat block until exp
	| if exp then block {elseif exp then block} [else block] end
	| for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
	| for namelist in explist do block end
	| function funcname funcbody
	| local function Name funcbody
	| local namelist [‘=’ explist]
	| varlist ‘=’ explist
	| functioncall
*/
func parseStat(lex *lexer.Lexer) ast.Stat {
	switch lex.LookAhead() {
	case lexer.TOKEN_SEP_SEMI:
		return parseEmptyStat(lex)
	case lexer.TOKEN_KW_BREAK:
		return parseBreakStat(lex)
	case lexer.TOKEN_SEP_LABEL:
		return parseLabelStat(lex)
	case lexer.TOKEN_KW_GOTO:
		return parseGotoStat(lex)
	case lexer.TOKEN_KW_DO:
		return parseDoStat(lex)
	case lexer.TOKEN_KW_WHILE:
		return parseWhileStat(lex)
	case lexer.TOKEN_KW_REPEAT:
		return parseRepeatStat(lex)
	case lexer.TOKEN_KW_IF:
		return parseIfStat(lex)
	case lexer.TOKEN_KW_FOR:
		return parseForStat(lex)
	case lexer.TOKEN_KW_FUNCTION:
		return parseFuncDefStat(lex)
	case lexer.TOKEN_KW_LOCAL:
		return parseLocalAssignOrFuncDefStat(lex)
	default:
		return parseAssignOrFuncCallStat(lex)
	}
}

// ;
func parseEmptyStat(lex *lexer.Lexer) *ast.EmptyStat {
	lex.NextTokenOfKind(lexer.TOKEN_SEP_SEMI)
	return _statEmpty
}

// break
func parseBreakStat(lex *lexer.Lexer) *ast.BreakStat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_BREAK)
	return &ast.BreakStat{lex.Line()}
}

// ‘::’ Name ‘::’
func parseLabelStat(lex *lexer.Lexer) *ast.LabelStat {
	lex.NextTokenOfKind(lexer.TOKEN_SEP_LABEL) // ::
	_, name := lex.NextIdentifier()            // name
	lex.NextTokenOfKind(lexer.TOKEN_SEP_LABEL) // ::
	return &ast.LabelStat{name}
}

// goto Name
func parseGotoStat(lex *lexer.Lexer) *ast.GotoStat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_GOTO) // goto
	_, name := lex.NextIdentifier()          // name
	return &ast.GotoStat{name}
}

// do block end
func parseDoStat(lex *lexer.Lexer) *ast.DoStat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_DO)  // do
	block := parseBlock(lex)                // block
	lex.NextTokenOfKind(lexer.TOKEN_KW_END) // end
	return &ast.DoStat{block}
}

// while exp do block end
func parseWhileStat(lex *lexer.Lexer) *ast.WhileStat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_WHILE) // while
	exp := parseExp(lex)                      // exp
	lex.NextTokenOfKind(lexer.TOKEN_KW_DO)    // do
	block := parseBlock(lex)                  // block
	lex.NextTokenOfKind(lexer.TOKEN_KW_END)   // end
	return &ast.WhileStat{exp, block}
}

// repeat block until exp
func parseRepeatStat(lex *lexer.Lexer) *ast.RepeatStat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_REPEAT) // repeat
	block := parseBlock(lex)                   // block
	lex.NextTokenOfKind(lexer.TOKEN_KW_UNTIL)  // until
	exp := parseExp(lex)                       // exp
	return &ast.RepeatStat{block, exp}
}

// if exp then block {elseif exp then block} [else block] end
func parseIfStat(lex *lexer.Lexer) *ast.IfStat {
	exps := make([]ast.Exp, 0, 4)
	blocks := make([]*ast.Block, 0, 4)

	lex.NextTokenOfKind(lexer.TOKEN_KW_IF)   // if
	exps = append(exps, parseExp(lex))       // exp
	lex.NextTokenOfKind(lexer.TOKEN_KW_THEN) // then
	blocks = append(blocks, parseBlock(lex)) // block

	for lex.LookAhead() == lexer.TOKEN_KW_ELSEIF {
		lex.NextToken()                          // elseif
		exps = append(exps, parseExp(lex))       // exp
		lex.NextTokenOfKind(lexer.TOKEN_KW_THEN) // then
		blocks = append(blocks, parseBlock(lex)) // block
	}

	// else block => elseif true then block
	if lex.LookAhead() == lexer.TOKEN_KW_ELSE {
		lex.NextToken()                               // else
		exps = append(exps, &ast.TrueExp{lex.Line()}) //
		blocks = append(blocks, parseBlock(lex))      // block
	}

	lex.NextTokenOfKind(lexer.TOKEN_KW_END) // end
	return &ast.IfStat{exps, blocks}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
// for namelist in explist do block end
func parseForStat(lex *lexer.Lexer) ast.Stat {
	lineOfFor, _ := lex.NextTokenOfKind(lexer.TOKEN_KW_FOR)
	_, name := lex.NextIdentifier()
	if lex.LookAhead() == lexer.TOKEN_OP_ASSIGN {
		return _finishForNumStat(lex, lineOfFor, name)
	} else {
		return _finishForInStat(lex, name)
	}
}

// for Name ‘=’ exp ‘,’ exp [‘,’ exp] do block end
func _finishForNumStat(lex *lexer.Lexer, lineOfFor int, varName string) *ast.ForNumStat {
	lex.NextTokenOfKind(lexer.TOKEN_OP_ASSIGN) // for name =
	initExp := parseExp(lex)                   // exp
	lex.NextTokenOfKind(lexer.TOKEN_SEP_COMMA) // ,
	limitExp := parseExp(lex)                  // exp

	var stepExp ast.Exp
	if lex.LookAhead() == lexer.TOKEN_SEP_COMMA {
		lex.NextToken()         // ,
		stepExp = parseExp(lex) // exp
	} else {
		stepExp = &ast.IntegerExp{lex.Line(), 1}
	}

	lineOfDo, _ := lex.NextTokenOfKind(lexer.TOKEN_KW_DO) // do
	block := parseBlock(lex)                              // block
	lex.NextTokenOfKind(lexer.TOKEN_KW_END)               // end

	return &ast.ForNumStat{lineOfFor, lineOfDo,
		varName, initExp, limitExp, stepExp, block}
}

// for namelist in explist do block end
// namelist ::= Name {‘,’ Name}
// explist ::= exp {‘,’ exp}
func _finishForInStat(lex *lexer.Lexer, name0 string) *ast.ForInStat {
	nameList := _finishNameList(lex, name0)               // for namelist
	lex.NextTokenOfKind(lexer.TOKEN_KW_IN)                // in
	expList := parseExpList(lex)                          // explist
	lineOfDo, _ := lex.NextTokenOfKind(lexer.TOKEN_KW_DO) // do
	block := parseBlock(lex)                              // block
	lex.NextTokenOfKind(lexer.TOKEN_KW_END)               // end
	return &ast.ForInStat{lineOfDo, nameList, expList, block}
}

// namelist ::= Name {‘,’ Name}
func _finishNameList(lex *lexer.Lexer, name0 string) []string {
	names := []string{name0}
	for lex.LookAhead() == lexer.TOKEN_SEP_COMMA {
		lex.NextToken()                 // ,
		_, name := lex.NextIdentifier() // Name
		names = append(names, name)
	}
	return names
}

// local function Name funcbody
// local namelist [‘=’ explist]
func parseLocalAssignOrFuncDefStat(lex *lexer.Lexer) ast.Stat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_LOCAL)
	if lex.LookAhead() == lexer.TOKEN_KW_FUNCTION {
		return _finishLocalFuncDefStat(lex)
	} else {
		return _finishLocalVarDeclStat(lex)
	}
}

/*
http://www.lua.org/manual/5.3/manual.html#3.4.11

function f() end          =>  f = function() end
function t.a.b.c.f() end  =>  t.a.b.c.f = function() end
function t.a.b.c:f() end  =>  t.a.b.c.f = function(self) end
local function f() end    =>  local f; f = function() end

The statement `local function f () body end`
translates to `local f; f = function () body end`
not to `local f = function () body end`
(This only makes a difference when the body of the function
 contains references to f.)
*/
// local function Name funcbody
func _finishLocalFuncDefStat(lex *lexer.Lexer) *ast.LocalFuncDefStat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_FUNCTION) // local function
	_, name := lex.NextIdentifier()              // name
	fdExp := parseFuncDefExp(lex)                // funcbody
	return &ast.LocalFuncDefStat{name, fdExp}
}

// local namelist [‘=’ explist]
func _finishLocalVarDeclStat(lex *lexer.Lexer) *ast.LocalVarDeclStat {
	_, name0 := lex.NextIdentifier()        // local Name
	nameList := _finishNameList(lex, name0) // { , Name }
	var expList []ast.Exp = nil
	if lex.LookAhead() == lexer.TOKEN_OP_ASSIGN {
		lex.NextToken()             // ==
		expList = parseExpList(lex) // explist
	}
	lastLine := lex.Line()
	return &ast.LocalVarDeclStat{lastLine, nameList, expList}
}

// varlist ‘=’ explist
// functioncall
func parseAssignOrFuncCallStat(lex *lexer.Lexer) ast.Stat {
	prefixExp := parsePrefixExp(lex)
	if fc, ok := prefixExp.(*ast.FuncCallExp); ok {
		return fc
	} else {
		return parseAssignStat(lex, prefixExp)
	}
}

// varlist ‘=’ explist |
func parseAssignStat(lex *lexer.Lexer, var0 ast.Exp) *ast.AssignStat {
	varList := _finishVarList(lex, var0)       // varlist
	lex.NextTokenOfKind(lexer.TOKEN_OP_ASSIGN) // =
	expList := parseExpList(lex)               // explist
	lastLine := lex.Line()
	return &ast.AssignStat{lastLine, varList, expList}
}

// varlist ::= var {‘,’ var}
func _finishVarList(lex *lexer.Lexer, var0 ast.Exp) []ast.Exp {
	vars := []ast.Exp{_checkVar(lex, var0)}        // var
	for lex.LookAhead() == lexer.TOKEN_SEP_COMMA { // {
		lex.NextToken()                          // ,
		exp := parsePrefixExp(lex)               // var
		vars = append(vars, _checkVar(lex, exp)) //
	} // }
	return vars
}

// var ::=  Name | prefixexp ‘[’ exp ‘]’ | prefixexp ‘.’ Name
func _checkVar(lex *lexer.Lexer, exp ast.Exp) ast.Exp {
	switch exp.(type) {
	case *ast.NameExp, *ast.TableAccessExp:
		return exp
	}
	lex.NextTokenOfKind(-1) // trigger error
	panic("unreachable!")
}

// function funcname funcbody
// funcname ::= Name {‘.’ Name} [‘:’ Name]
// funcbody ::= ‘(’ [parlist] ‘)’ block end
// parlist ::= namelist [‘,’ ‘...’] | ‘...’
// namelist ::= Name {‘,’ Name}
func parseFuncDefStat(lex *lexer.Lexer) *ast.AssignStat {
	lex.NextTokenOfKind(lexer.TOKEN_KW_FUNCTION) // function
	fnExp, hasColon := _parseFuncName(lex)       // funcname
	fdExp := parseFuncDefExp(lex)                // funcbody
	if hasColon {                                // insert self
		fdExp.ParList = append(fdExp.ParList, "")
		copy(fdExp.ParList[1:], fdExp.ParList)
		fdExp.ParList[0] = "self"
	}

	return &ast.AssignStat{
		LastLine: fdExp.Line,
		VarList:  []ast.Exp{fnExp},
		ExpList:  []ast.Exp{fdExp},
	}
}

// funcname ::= Name {‘.’ Name} [‘:’ Name]
func _parseFuncName(lex *lexer.Lexer) (exp ast.Exp, hasColon bool) {
	line, name := lex.NextIdentifier()
	exp = &ast.NameExp{line, name}

	for lex.LookAhead() == lexer.TOKEN_SEP_DOT {
		lex.NextToken()
		line, name := lex.NextIdentifier()
		idx := &ast.StringExp{line, name}
		exp = &ast.TableAccessExp{line, exp, idx}
	}
	if lex.LookAhead() == lexer.TOKEN_SEP_COLON {
		lex.NextToken()
		line, name := lex.NextIdentifier()
		idx := &ast.StringExp{line, name}
		exp = &ast.TableAccessExp{line, exp, idx}
		hasColon = true
	}

	return
}
