package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	//"github.com/alecthomas/repr"
	"log"
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
	CmpGR
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

// RangedReplace 
// startIndex and endIndex are inclusive
func (expressionIntermediates *ExpressionIntermediates) RangedReplace(startIndex int, endIndex int, replaceExpressionIntermediate ExpressionIntermediate) {
	(*expressionIntermediates)[startIndex] = replaceExpressionIntermediate
	*expressionIntermediates = append((*expressionIntermediates)[:startIndex+1], (*expressionIntermediates)[endIndex+1:]...)
}

func (expressionIntermediates *ExpressionIntermediates) Evaluate(scope *Scope) value.Value {
	// Shorten Intermediates and generate LLVM IR
	// In the main loop, do in thing only at a time.
	// The whole array could be modified in one branch of the loop.
	for {
		// Deal with function calls

		// Split (...) into ExpressionIntermediates to solve

		// Solve first "*" / "/"

		// Solve first "+" / "-"
		for index, expressionIntermediate := range *expressionIntermediates {
			if expressionIntermediate.OperationType == Add {
				if expressionIntermediates.ValidateTwoElementOperation(index) == false {
					log.Fatal("Unable to generate IR for + Operation")
				}
				tmp := scope.Block.NewAdd((*expressionIntermediates)[index-1].Value, (*expressionIntermediates)[index+1].Value)
				tmpExpressionIntermediate := NewValue(tmp)
				expressionIntermediates.RangedReplace(index-1, index+1, tmpExpressionIntermediate)
				goto finish
			}
			if expressionIntermediate.OperationType == Minus {
				if expressionIntermediates.ValidateTwoElementOperation(index) == false {
					log.Fatal("Unable to generate IR for - Operation")
				}
				tmp := scope.Block.NewSub((*expressionIntermediates)[index-1].Value, (*expressionIntermediates)[index+1].Value)
				tmpExpressionIntermediate := NewValue(tmp)
				expressionIntermediates.RangedReplace(index-1, index+1, tmpExpressionIntermediate)
				goto finish
			}
		}

finish:
		// Exit
		//log.Println(len(expressionIntermediates))
		if len(*expressionIntermediates) == 1 {
			return (*expressionIntermediates)[0].Value
		}
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
	case token.Add != nil:
		return Add
	case token.Minus != nil:
		return Minus
	case token.Times != nil:
		return Times
	case token.Divide != nil:
		return Divide
	case token.CmpEQ != nil:
		return CmpEQ
	case token.CmpNE != nil:
		return CmpNE
	case token.CmpGR != nil:
		return CmpGR
	case token.CmpGE != nil:
		return CmpGE
	case token.CmpLT != nil:
		return CmpLT
	case token.CmpLE != nil:
		return CmpLE
	case token.LeftParenthesis != nil:
		return LeftParenthesis
	case token.RightParenthesis != nil:
		return RightParenthesis
	case token.LeftBracket != nil:
		return LeftBracket
	case token.RightBracket != nil:
		return RightBracket
	}
	return NotOperationSymbol
}

func (c *Constant) Eval(scope *Scope) value.Value {
	switch {
	case c.VString != nil:
		// Static Strings should be stored globally.
		string_constant := constant.NewCharArray([]byte(*c.VString))
		global_def := scope.Module.NewGlobalDef("", string_constant)
		return global_def
	case c.VReal != nil:
		return constant.NewFloat(types.Float, *c.VReal)
	case c.VInt != nil:
		return constant.NewInt(types.I32, *c.VInt)
	default:
		log.Fatal("Parser: cannot parse value at", c.Pos)
	}
	return nil
}
