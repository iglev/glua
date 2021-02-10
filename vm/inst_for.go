package vm

import (
	"fmt"

	"github.com/iglev/glua/api"
)

/*
lua code:
	for i=1,10,1 do
		-- nothing
	end
pc   op       stack:
[01] LOADK    [1][nil][nil][nil]
[02] LOADK    [1][10][nil][nil]
[03] LOADK    [1][10][2][nil]
[04] FORPREP  [-1][10][2][nil]
[05] FORLOOP  [1][10][2][1]
[05] FORLOOP  [3][10][2][3]
[05] FORLOOP  [5][10][2][5]
[05] FORLOOP  [7][10][2][7]
[05] FORLOOP  [9][10][2][9]
[05] FORLOOP  [11][10][2][9]
*/

// R(A)-=R(A+2); pc+=sBx
func forPrep(i Instruction, vm api.LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	fmt.Printf("t1=%v t2=%v t3=%v vm=%+v\n", vm.Type(a), vm.Type(a+1), vm.Type(a+2), vm)

	if vm.Type(a) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)
	}
	if vm.Type(a+1) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)
	}
	if vm.Type(a+2) == api.LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)
	}

	vm.PushValue(a)
	vm.PushValue(a + 2)
	vm.Arith(api.LUA_OPSUB)
	vm.Replace(a)
	vm.AddPC(sBx)
}

// R(A)+=R(A+2);
// if R(A) <?= R(A+1) then {
//   pc+=sBx; R(A+3)=R(A)
// }
func forLoop(i Instruction, vm api.LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	// R(A)+=R(A+2);
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(api.LUA_OPADD)
	vm.Replace(a)

	isPositiveStep := vm.ToNumber(a+2) >= 0
	if isPositiveStep && vm.Compare(a, a+1, api.LUA_OPLE) ||
		!isPositiveStep && vm.Compare(a+1, a, api.LUA_OPLE) {
		// pc+=sBx; R(A+3)=R(A)
		vm.AddPC(sBx)
		vm.Copy(a, a+3)
	}
}
