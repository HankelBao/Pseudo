package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	//"github.com/alecthomas/repr"
	"log"
	"strconv"
)

type OperationSymbol int

const (
	NotOperationSymbol OperationSymbol = iota
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

func (e *Expression) Evaluate(scope *Scope) value.Value {
	// Convert all the ExpressionToken into ExpressionIntermediate
	expressionIntermediates := make(ExpressionIntermediates, len(e.Tokens))
	for index, token := range e.Tokens {
		operationType := token.GetOperationSymbol()
		// Operation
		if operationType != NotOperationSymbol {
			expressionIntermediates[index] = NewOperation(operationType)
			continue
		}
		// Constant
		if token.Constant != nil {
			expressionIntermediates[index] = NewValue(token.Constant.Eval(scope))
			continue
		}
		expressionIntermediates[index] = NewNone()
	}
	for index, token := range e.Tokens {
		// Symbols
		if token.Symbol != nil {
			if index != len(e.Tokens)-1 {
				if expressionIntermediates[index+1].OperationType == LeftParenthesis {
					// Function
					function := scope.FindFunction(*token.Symbol)
					expressionIntermediates[index] = NewFunc(function)
					continue
				}
			}
			// Variable
			variablePtr := scope.FindVariable(*token.Symbol)
			variable := scope.Block.NewLoad(variablePtr)
			expressionIntermediates[index] = NewValue(variable)
		}
	}

	firstOT := expressionIntermediates[0].OperationType
	if firstOT == Add || firstOT == Minus || firstOT == Times || firstOT == Divide {
		zero := constant.NewInt(types.I32, 0)
		zeroExpressionIntermediate := NewValue(zero)
		expressionIntermediates = append([]ExpressionIntermediate{zeroExpressionIntermediate}, expressionIntermediates...)
	}

	return (&expressionIntermediates).Evaluate(scope)
}

type ExpressionIntermediates []ExpressionIntermediate

func (expressionIntermediates *ExpressionIntermediates) ValidateTwoElementOperation(index int) bool {
	if index < 1 || index > len(*expressionIntermediates)-2 {
		return false
	}
	if (*expressionIntermediates)[index-1].Value == nil {
		return false
	}
	if (*expressionIntermediates)[index+1].Value == nil {
		return false
	}
	return true
}

func (expressionIntermediates *ExpressionIntermediates) TwoElemOperation(index int, opFunc func(value.Value, value.Value) value.Value) {
	if expressionIntermediates.ValidateTwoElementOperation(index) == false {
		log.Fatal("Unable to generate IR for Two Element Operation")
	}
	value1 := (*expressionIntermediates)[index-1].Value
	value2 := (*expressionIntermediates)[index+1].Value
	result := opFunc(value1, value2)
	tmpExpressionIntermediate := NewValue(result)
	expressionIntermediates.RangedReplace(index-1, index+1, tmpExpressionIntermediate)
}

// RangedReplace
// startIndex and endIndex are inclusive
func (expressionIntermediates *ExpressionIntermediates) RangedReplace(startIndex int, endIndex int, replaceExpressionIntermediate ExpressionIntermediate) {
	(*expressionIntermediates)[startIndex] = replaceExpressionIntermediate
	*expressionIntermediates = append((*expressionIntermediates)[:startIndex+1], (*expressionIntermediates)[endIndex+1:]...)
}

func (expressionIntermediates *ExpressionIntermediates) Evaluate(scope *Scope) value.Value {
	// Shorten Intermediates and generate LLVM IR
	// In the main loop, do one thing only at a time.
	// The whole array could be modified in one branch of the loop.
	for {
		if len(*expressionIntermediates) == 1 {
			return (*expressionIntermediates)[0].Value
		}

		// Deal with function calls

		// Split (...) into ExpressionIntermediates to solve
		// Now all the "(" and ")" should be prioritized expressions
		lastLeftParenthesisIndex := -1
		for index, expressionIntermediate := range *expressionIntermediates {
			if expressionIntermediate.OperationType == LeftParenthesis {
				lastLeftParenthesisIndex = index
			}
			if expressionIntermediate.OperationType == RightParenthesis {
				if lastLeftParenthesisIndex == -1 {
					log.Fatal("Unmatched number of parenthesis.")
				}
				subExpressionIntermediates := (*expressionIntermediates)[lastLeftParenthesisIndex+1 : index]
				log.Printf("%+v", subExpressionIntermediates)
				result := subExpressionIntermediates.Evaluate(scope)
				resultExpressionIntermediate := NewValue(result)
				expressionIntermediates.RangedReplace(lastLeftParenthesisIndex, index, resultExpressionIntermediate)
				goto end
			}
		}

		// Solve first "*" / "/"

		// Solve first "+" / "-"
		for index, expressionIntermediate := range *expressionIntermediates {
			if expressionIntermediate.OperationType == Add {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewAdd(value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFAdd(value1, value2)
					}
					return nil
				})
				goto end
			}
			if expressionIntermediate.OperationType == Minus {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewSub(value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFSub(value1, value2)
					}
					return nil
				})
				goto end
			}
		}

		// CmpOp
		for index, expressionIntermediate := range *expressionIntermediates {
			if expressionIntermediate.OperationType == CmpEQ {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewICmp(enum.IPredEQ, value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFCmp(enum.FPredOEQ, value1, value2)
					}
					return nil
				})
				goto end
			}
			if expressionIntermediate.OperationType == CmpNE {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewICmp(enum.IPredNE, value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFCmp(enum.FPredONE, value1, value2)
					}
					return nil
				})
				goto end
			}
			if expressionIntermediate.OperationType == CmpLT {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewICmp(enum.IPredSLT, value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFCmp(enum.FPredOLT, value1, value2)
					}
					return nil
				})
				goto end
			}
			if expressionIntermediate.OperationType == CmpLE {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewICmp(enum.IPredSLE, value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFCmp(enum.FPredOLE, value1, value2)
					}
					return nil
				})
				goto end
			}
			if expressionIntermediate.OperationType == CmpGT {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewICmp(enum.IPredSGT, value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFCmp(enum.FPredOGT, value1, value2)
					}
					return nil
				})
				goto end
			}
			if expressionIntermediate.OperationType == CmpGE {
				expressionIntermediates.TwoElemOperation(index, func(value1 value.Value, value2 value.Value) value.Value {
					if value1.Type() == types.I32 {
						return scope.Block.NewICmp(enum.IPredSGE, value1, value2)
					} else if value1.Type() == types.Double {
						return scope.Block.NewFCmp(enum.FPredOGE, value1, value2)
					}
					return nil
				})
				goto end
			}
		}

		// If nothing was done in a loop, the expression is not solvable.
		// Panic.
		// log.Fatal("Expression of unsolvable sequence.")
	end:
	}
}

type ExpressionIntermediate struct {
	OperationType OperationSymbol
	Value         value.Value
	Func          *ir.Func
}

func NewNone() ExpressionIntermediate {
	return ExpressionIntermediate{
		OperationType: NotOperationSymbol,
		Value:         nil,
		Func:          nil,
	}
}

func NewOperation(operationType OperationSymbol) ExpressionIntermediate {
	return ExpressionIntermediate{
		OperationType: operationType,
		Value:         nil,
		Func:          nil,
	}
}

func NewValue(value value.Value) ExpressionIntermediate {
	return ExpressionIntermediate{
		OperationType: NotOperationSymbol,
		Value:         value,
		Func:          nil,
	}
}

func NewFunc(function *ir.Func) ExpressionIntermediate {
	return ExpressionIntermediate{
		OperationType: NotOperationSymbol,
		Value:         nil,
		Func:          function,
	}
}

func (token *ExpressionToken) GetOperationSymbol() OperationSymbol {
	switch {
	case token.BasicOp != nil:
		switch *(token.BasicOp) {
		case "+":
			return Add
		case "-":
			return Minus
		case "*":
			return Times
		case "/":
			return Divide
		}
	case token.CmpOp != nil:
		switch *(token.CmpOp) {
		case "=":
			return CmpEQ
		case "!=":
			return CmpNE
		case ">":
			return CmpGT
		case ">=":
			return CmpGE
		case "<":
			return CmpLT
		case "<=":
			return CmpLE
		}
	case token.Parenthesis != nil:
		switch *(token.Parenthesis) {
		case "(":
			return LeftParenthesis
		case ")":
			return RightParenthesis
		}
	case token.Bracket != nil:
		switch *(token.Bracket) {
		case "[":
			return LeftBracket
		case "]":
			return RightBracket
		}
	}
	return NotOperationSymbol
}

var stringConstantNameIndex = 0

func (c *Constant) Eval(scope *Scope) value.Value {
	switch {
	case c.VString != nil:
		// Static Strings should be stored globally.
		stringConstant := constant.NewCharArrayFromString(*c.VString + "\000")
		stringGName := "goPseConstant?$" + strconv.Itoa(stringConstantNameIndex)
		stringConstantNameIndex++
		globalDef := scope.Module.NewGlobalDef(stringGName, stringConstant)
		return globalDef
	case c.VReal != nil:
		return constant.NewFloat(types.Double, *c.VReal)
	case c.VInt != nil:
		return constant.NewInt(types.I32, *c.VInt)
	default:
		log.Fatal("Parser: cannot parse value at", c.Pos)
	}
	return nil
}
