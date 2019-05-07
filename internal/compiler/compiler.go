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

	main_scope := globalScope.NewFuncScope(main)
	main_scope.Main = true
	ast.Compile(main_scope)

	zero := constant.NewInt(types.I32, 0)
	main_scope.Block.NewRet(zero)
	return globalScope.Module
}
