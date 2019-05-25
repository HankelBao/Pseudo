package compiler

import (
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

// InitRuntime inits the runtime for pseudocode.
// This includes:
// - C Standard Functions
// - Format for PrintfD and PrintfF
func (scope *Scope) InitRuntime() {
	if scope.IsGlobal() == false {
		log.Fatal("Init C Shared Lib to non-global scope")
	}
	mod := scope.Module

	puts := mod.NewFunc("puts", &types.IntType{BitSize: 32}, ir.NewParam("", &types.PointerType{ElemType: types.I8}))
	scope.RegisterFunction("puts", puts)

	printf := mod.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
	printf.Sig.Variadic = true
	scope.RegisterFunction("printf", printf)

	printfdFmt := constant.NewCharArrayFromString("Int: %d\n\000")
	printfdFmtDef := mod.NewGlobalDef("printfd_fmt", printfdFmt)
	scope.RegisterVariable("printfd_fmt", printfdFmtDef)

	printffFmt := constant.NewCharArrayFromString("Int: %f\n\000")
	printffFmtDef := mod.NewGlobalDef("printff_fmt", printffFmt)
	scope.RegisterVariable("printff_fmt", printffFmtDef)
}
