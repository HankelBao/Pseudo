package compiler

import (
	"github.com/llir/llvm/ir/value"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/enum"
)

// OperationType describes the type of an operation.
type OperationType int

const (
	// NotOperationType describes the filed as an none-operation.
	NotOperationType OperationType = iota
	Add
	Minus
	Times
	Divide
	CmpEQ
	CmpNE
	CmpGT
	CmpGE
	CmpLT
	CmpLE
	LeftParenthesis
	RightParenthesis
	LeftBracket
	RightBracket
	ParamDivier
)

// EvaluateFunc is type defining the operation evaluation function.
type EvaluateFunc func(*Scope, value.Value, value.Value) value.Value

// OperationPriority stores the priorities of the operations.
// Not all the operations of OperationType have a priority.
//
// Priority is the priority to deal with operations.
// The higher the priority is, the earlier the operation should be dealt with.
// As the result, operation with lower priority would be assigned closer to the root of the tree.
// Operations would be dealt with starting from leaves.
var OperationPriority = map[OperationType]int {
	Add: 5,
	Minus: 5,
	Times: 6,
	Divide: 6,
	CmpEQ: 2,
	CmpNE: 2,
	CmpGT: 2,
	CmpGE: 2,
	CmpLT: 2,
	CmpLE: 2,
	// Parenthesis brings anything in parenthesises to higher unit.
	// 10 is the unit holding priorities of all the oprations.
	LeftParenthesis: 10,
	RightParenthesis: 10,
}

// OperationEvaluateFunc stores the evalutation function of the operations.
var OperationEvaluateFunc = map[OperationType]EvaluateFunc {
	Add: AddEval,
	Minus: MinusEval,
	CmpEQ: CmpEQEval,
	CmpNE: CmpNEEval,
	CmpGT: CmpGTEval,
	CmpGE: CmpGEEval,
	CmpLT: CmpLTEval,
	CmpLE: CmpLEEval,
}

// AddEval generates IR for add
func AddEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewAdd(value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFAdd(value1, value2)
	}
	return nil
}

// MinusEval generates IR for minus
func MinusEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewSub(value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFSub(value1, value2)
	}
	return nil
}

func CmpEQEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredEQ, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOEQ, value1, value2)
	}
	return nil
}

func CmpNEEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredNE, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredONE, value1, value2)
	}
	return nil
}

func CmpLTEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSLT, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOLT, value1, value2)
	}
	return nil
}

func CmpLEEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSLE, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOLE, value1, value2)
	}
	return nil
}

func CmpGTEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSGT, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOGT, value1, value2)
	}
	return nil
}

func CmpGEEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSGE, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOGE, value1, value2)
	}
	return nil
}