package compiler

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/constant"
	"log"
)

func (ast *Ast) Compile(scope *Scope) {
	for _, inst := range ast.Instructions {
		switch {
		case inst.Output != nil:
			inst.Output.Compile(scope)
		case inst.DeclareVariable != nil:
			inst.DeclareVariable.Compile(scope)
		case inst.Assignment != nil:
			inst.Assignment.Compile(scope)
		case inst.PrintfS != nil:
			inst.PrintfS.Compile(scope)
		case inst.NullLine != nil:
			continue
		default:
			log.Fatal("Unknown instruction")
		}
	}
}

func (ins *InstOutput) Compile(scope *Scope) {
	value := ins.Content.Evaluate(scope)
	tmp_ptr := scope.Block.NewBitCast(value, &types.PointerType{ElemType: &types.IntType{BitSize: 8}})
	puts := scope.FindFunction("puts")
	if puts == nil {
		log.Fatal("puts not found")
	}
	scope.Block.NewCall(puts, tmp_ptr)
}

func (ins *InstDeclareVariable) Compile(scope *Scope) {
	var variableName string = ins.Name
	var variableInitial constant.Constant
	switch {
	case ins.Type.Int != nil:
		variableInitial = constant.NewInt(types.I32, 0)
		break
	case ins.Type.REAL != nil:
		variableInitial = constant.NewFloat(types.Float, 0.0)
		break
	}

	// If the scope in the main scope,
	// Variables are defined globally.
	if scope.Main {
		newVariable := scope.Module.NewGlobalDef(variableName, variableInitial)
		scope.GlobalScope.RegisterVariable(variableName, newVariable)
	} else {
		// TODO: Private variable, allocate...
	}
}

func (ins *InstAssignment) Compile(scope *Scope) {
	key := ins.Left.Locate(scope)
	expression := ins.Right.Evaluate(scope)
	scope.Block.NewStore(expression, key)
}

func (ins *InstPrintfS) Compile(scope *Scope) {
	value := ins.Content.Evaluate(scope)

	format_def := scope.FindVariable("printfs_fmt")
	if format_def == nil {
		log.Fatal("Cannot find printfs format")
	}
	format_ptr := scope.Block.NewBitCast(format_def, types.I8Ptr)

	printf := scope.FindFunction("printf")
	scope.Block.NewCall(printf, format_ptr, value)
}
