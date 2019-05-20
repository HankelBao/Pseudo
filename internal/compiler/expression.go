package compiler

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	//"github.com/alecthomas/repr"
	"log"
	"strconv"
)

// Evaluate evalues the expression and generates IR.
func (e *Expression) Evaluate(scope *Scope) value.Value {
	for index := range e.Tokens {
		e.Tokens[index].EvaluateOperationType()
	}
	e.Prioritize()

	expressionIR := NewExpressionIR(e.Tokens)
	expressionIR.Build(scope)
	return expressionIR.Evaluate(scope)
}

// Prioritize prioritizes the operations in its tokens and updated the priority in its own field.
// All the operations has to be evaluated first!
func (e *Expression) Prioritize() {
	currentPriority := 0
	for index, token := range e.Tokens {
		if token.OperationType != NotOperationType {
			// The priority of the current operation should affect how the operation is dealt with.
			// So update currentPriority first.
			currentPriority += OperationPriority[token.OperationType]
			e.Tokens[index].Priority = currentPriority
		}
	}
}

// EvaluateOperationType gets the type of each operation.
func (token *ExpressionToken) EvaluateOperationType() {
	token.OperationType = func() OperationType {
		if token.OperationSymbol == nil {
			return NotOperationType
		}
		switch (*token.OperationSymbol) {
		case "+":
			return Add
		case "-":
			return Minus
		case "*":
			return Times
		case "/":
			return Divide
		case "=":
			return CmpEQ
		case "<>":
			return CmpNE
		case ">":
			return CmpGT
		case ">=":
			return CmpGE
		case "<":
			return CmpLT
		case "<=":
			return CmpLE
		case "(":
			return LeftParenthesis
		case ")":
			return RightParenthesis
		case "[":
			return LeftBracket
		case "]":
			return RightBracket
		case ",":
			return ParamDivier
		default:
			return NotOperationType
		}
	}()
}

var stringConstantNameIndex = 0

// Evaluate gets the value of the constant.
// If it is a string, it would be stored as a global variable.
func (c *Constant) Evaluate(scope *Scope) value.Value {
	switch {
	case c.VString != nil:
		// Static Strings should be stored globally.
		stringConstant := constant.NewCharArrayFromString(*c.VString + "\000")
		stringGName := "PseudoConstant?$" + strconv.Itoa(stringConstantNameIndex)
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
