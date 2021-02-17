package codegen

import "github.com/iglev/glua/compiler/ast"

func isVarargOrFuncCall(exp ast.Exp) bool {
	switch exp.(type) {
	case *ast.VarargExp, *ast.FuncCallExp:
		return true
	}
	return false
}

func removeTailNils(exps []ast.Exp) []ast.Exp {
	for n := len(exps) - 1; n >= 0; n-- {
		if _, ok := exps[n].(*ast.NilExp); !ok {
			return exps[0 : n+1]
		}
	}
	return nil
}
