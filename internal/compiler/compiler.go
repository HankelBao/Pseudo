package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

func Compile(ast *Ast) *ir.Module {
	globalScope := NewGlobalScope()
	globalScope.InitRuntime()

	main := globalScope.Module.NewFunc("main", &types.IntType{BitSize: 32})
	globalScope.RegisterFunction("main", main)

	main_block := main.NewBlock("")
	main_scope := globalScope.NewScope(main_block)
	main_scope.Main = true
	ast.Compile(main_scope)

	zero := constant.NewInt(types.I32, 0)
	main_block.NewRet(zero)
	return globalScope.Module
}
