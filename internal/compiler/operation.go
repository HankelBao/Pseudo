package compiler

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

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

// MultipleEval generates IR for multiple
func MultipleEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewMul(value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFMul(value1, value2)
	}
	return nil
}

// DivideEval generates IR for divide
func DivideEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewSDiv(value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFDiv(value1, value2)
	}
	return nil
}

// OppositeEval generates IR for opposite
func OppositeEval(scope *Scope, value value.Value) value.Value {
	if value.Type() == types.I32 {
		return scope.Block.NewSub(constant.NewInt(types.I32, 0), value)
	} else if value.Type() == types.Double {
		return scope.Block.NewFSub(constant.NewFloat(types.Float, 0.0), value)
	}
	return nil
}

// CmpEQEval generates IR for =
func CmpEQEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredEQ, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOEQ, value1, value2)
	}
	return nil
}

// CmpNEEval generates IR for <>
func CmpNEEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredNE, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredONE, value1, value2)
	}
	return nil
}

// CmpLTEval generates IR for <
func CmpLTEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSLT, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOLT, value1, value2)
	}
	return nil
}

// CmpLEEval generates IR for <=
func CmpLEEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSLE, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOLE, value1, value2)
	}
	return nil
}

// CmpGTEval generates IR for >
func CmpGTEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSGT, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOGT, value1, value2)
	}
	return nil
}

// CmpGEEval generates IR for >=
func CmpGEEval(scope *Scope, value1 value.Value, value2 value.Value) value.Value {
	if value1.Type() == types.I32 {
		return scope.Block.NewICmp(enum.IPredSGE, value1, value2)
	} else if value1.Type() == types.Double {
		return scope.Block.NewFCmp(enum.FPredOGE, value1, value2)
	}
	return nil
}
