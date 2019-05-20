package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

// Compile compiles the ast.
func Compile(ast *Ast) *ir.Module {
	globalScope := NewGlobalScope()
	globalScope.InitRuntime()

	main := globalScope.Module.NewFunc("main", &types.IntType{BitSize: 32})
	globalScope.RegisterFunction("main", main)

	mainScope := globalScope.NewFuncScope(main)
	mainScope.Main = true
	ast.Compile(mainScope)

	zero := constant.NewInt(types.I32, 0)
	mainScope.Block.NewRet(zero)
	return globalScope.Module
}
