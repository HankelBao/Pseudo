package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"log"
)

func (scope *Scope) InitCSharedLib() {
	if scope.IsGlobal() == false {
		log.Fatal("Init C Shared Lib to non-global scope")
	}
	mod := scope.Module

	puts := mod.NewFunc("puts", &types.IntType{BitSize: 32}, ir.NewParam("", &types.PointerType{ElemType: types.I8}))
	scope.RegisterFunction("puts", puts)

	printf := mod.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
	printf.Sig.Variadic = true
	scope.RegisterFunction("printf", printf)

	printfs_fmt := constant.NewCharArrayFromString("%d")
	printfs_fmt_def := mod.NewGlobalDef("printfs_fmt", printfs_fmt)
	scope.RegisterVariable("printfs_fmt", printfs_fmt_def)
}
