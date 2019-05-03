package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func Compile(ast *Ast) *ir.Module {
	globalScope := NewGlobalScope()
	globalScope.InitCSharedLib()

	main := globalScope.Module.NewFunc("main", &types.IntType{BitSize: 32})
	globalScope.RegisterFunction("main", main)

	main_block := main.NewBlock("")
	main_scope := globalScope.NewScope(main_block)
	main_scope.Main = true
	ast.Compile(main_scope)

	/*printf := main_scope.FindFunction("printf")
	string_constant := constant.NewCharArrayFromString("Hello World %d times\n")
	times := constant.NewInt(types.I32, 8)
	globalDef := globalScope.Module.NewGlobalDef("", string_constant)
	tmp := main_block.NewBitCast(globalDef, types.I8Ptr)
	main_block.NewCall(printf, tmp, times)*/

	zero := constant.NewInt(types.I32, 0)
	main_block.NewRet(zero)
	return globalScope.Module
}
