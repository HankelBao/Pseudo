package compiler

import (
	"log"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

// Interpreter is the interface for all the interpreters for asts.
type Interpreter interface {
	Compile(*Scope)
}

// Compile compiles the ast
// It splits them into different instructions.
func (ast *Ast) Compile(scope *Scope) {
	for _, inst := range ast.Instructions {
		switch {
		case inst.Output != nil:
			inst.Output.Compile(scope)
		case inst.DeclareVariable != nil:
			inst.DeclareVariable.Compile(scope)
		case inst.Assignment != nil:
			inst.Assignment.Compile(scope)
		case inst.PrintfD != nil:
			inst.PrintfD.Compile(scope)
		case inst.PrintfF != nil:
			inst.PrintfF.Compile(scope)
		case inst.ConditionBr != nil:
			inst.ConditionBr.Compile(scope)
		case inst.NullLine != nil:
			continue
		default:
			log.Fatal("Unknown instruction")
		}
	}
}

// Compile compiles InstOutput
func (ins *InstOutput) Compile(scope *Scope) {
	value := ins.Content.Evaluate(scope)
	tmpPtr := scope.Block.NewBitCast(value, &types.PointerType{ElemType: &types.IntType{BitSize: 8}})
	puts := scope.FindFunction("puts")
	if puts == nil {
		log.Fatal("puts not found")
	}
	scope.Block.NewCall(puts, tmpPtr)
}

// Compile compiles InstDeclareVariable
func (ins *InstDeclareVariable) Compile(scope *Scope) {
	variableName := ins.Name
	var variableInitial constant.Constant
	switch {
	case ins.Type.Int != nil:
		variableInitial = constant.NewInt(types.I32, 0)
		break
	case ins.Type.REAL != nil:
		variableInitial = constant.NewFloat(types.Double, 0.0)
		break
	case ins.Type.BOOL != nil:
		variableInitial = constant.NewBool(false)
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

// Compile compiles InstAssignment
func (ins *InstAssignment) Compile(scope *Scope) {
	key := ins.Left.Locate(scope)
	expression := ins.Right.Evaluate(scope)
	scope.Block.NewStore(expression, key)
}

// Compile compiles InstPrintfD
func (ins *InstPrintfD) Compile(scope *Scope) {
	value := ins.Content.Evaluate(scope)

	formatDef := scope.FindVariable("printfd_fmt")
	if formatDef == nil {
		log.Fatal("Cannot find printfd format")
	}
	formatPtr := scope.Block.NewBitCast(formatDef, types.I8Ptr)

	printf := scope.FindFunction("printf")
	scope.Block.NewCall(printf, formatPtr, value)
}

// Compile compiles InstPrintfF
func (ins *InstPrintfF) Compile(scope *Scope) {
	value := ins.Content.Evaluate(scope)

	formatDef := scope.FindVariable("printff_fmt")
	if formatDef == nil {
		log.Fatal("Cannot find printff format")
	}
	formatPtr := scope.Block.NewBitCast(formatDef, types.I8Ptr)

	printf := scope.FindFunction("printf")
	scope.Block.NewCall(printf, formatPtr, value)

}

// Compile compiles InstConditionBr
func (ins *InstConditionBr) Compile(scope *Scope) {
	condVal := ins.Condition.Evaluate(scope)

	trueBlock := scope.Func.NewBlock("")
	trueBlockScope := scope.NewScope(trueBlock)
	ins.TrueBr.Compile(trueBlockScope)

	falseBlock := scope.Func.NewBlock("")
	falseBlockScope := scope.NewScope(falseBlock)
	if ins.FalseBr != nil {
		ins.FalseBr.Compile(falseBlockScope)
	}

	continueBlock := scope.Func.NewBlock("")
	trueBlock.NewBr(continueBlock)
	falseBlock.NewBr(continueBlock)

	scope.Block.NewCondBr(condVal, trueBlock, falseBlock)
	scope.Block = continueBlock
}
