package compiler

import (
	"log"
	"math"

	"github.com/llir/llvm/ir/value"
)

// ExpressionIR converts ExpressionToken into real asts(binary tree).
type ExpressionIR struct {
	OperationType    OperationType
	ExpressionTokens []ExpressionToken
	Value            value.Value
	Left             *ExpressionIR
	Right            *ExpressionIR
}

// NewExpressionIR creates a new expressionIR with initial expression tokens.
func NewExpressionIR(expressionTokens []ExpressionToken) *ExpressionIR {
	expressionIR := ExpressionIR{
		OperationType:    NotOperationType,
		ExpressionTokens: expressionTokens,
		Value:            nil,
		Left:             nil,
		Right:            nil,
	}
	return &expressionIR
}

// Build builds the binary tree from its ExpressionTokens.
func (expressionIR *ExpressionIR) Build(scope *Scope) {
	// Find the operation with least priority.
	targetPriority := math.MaxInt32
	targetIndex := -1
	for index := len(expressionIR.ExpressionTokens) - 1; index >= 0; index-- {
		token := expressionIR.ExpressionTokens[index]
		if token.OperationType == NotOperationType {
			continue
		}
		if token.Priority < targetPriority {
			if token.OperationType == LeftParenthesis || token.OperationType == RightParenthesis {
				continue
			}
			targetIndex = index
			targetPriority = token.Priority
		}
	}
	log.Println(targetIndex)

	// Check whether it is a leaf.
	var leafNode bool
	if targetIndex == -1 {
		leafNode = true
	} else {
		leafNode = false
	}

	if leafNode == false {
		expressionIR.OperationType = expressionIR.ExpressionTokens[targetIndex].OperationType
		expressionIR.Left = NewExpressionIR(expressionIR.ExpressionTokens[:targetIndex])
		expressionIR.Right = NewExpressionIR(expressionIR.ExpressionTokens[targetIndex+1:])
		expressionIR.Left.Build(scope)
		expressionIR.Right.Build(scope)
	} else {
		// Parenthesis could be mixed in the tokens here.
		// Parenthesis at the head would always be parenthesises for prioritizing expressions.
		numberHeadParenToRemove := 0
		for _, token := range expressionIR.ExpressionTokens {
			if token.OperationType == LeftParenthesis {
				numberHeadParenToRemove++
			} else {
				break
			}
		}
		expressionIR.ExpressionTokens = expressionIR.ExpressionTokens[numberHeadParenToRemove:]

		// Count the left parenthesis left, they are used for functions
		numberFuncLeftParen := 0
		for _, token := range expressionIR.ExpressionTokens {
			if token.OperationType == LeftParenthesis {
				numberFuncLeftParen++
			}
		}
		numberTailRightParen := 0
		for i := len(expressionIR.ExpressionTokens) - 1; i > 0; i-- {
			if expressionIR.ExpressionTokens[i].OperationType == RightParenthesis {
				numberTailRightParen++
			} else {
				break
			}
		}

		numberTailRightParenToRemove := numberTailRightParen - numberFuncLeftParen
		if numberTailRightParenToRemove < 0 {
			numberTailRightParenToRemove = 0
		}
		expressionIR.ExpressionTokens = expressionIR.ExpressionTokens[:len(expressionIR.ExpressionTokens)-numberTailRightParenToRemove]

		token := expressionIR.ExpressionTokens[0]
		if token.Constant != nil {
			expressionIR.Value = token.Constant.Evaluate(scope)
		}
		if token.Symbol != nil {
			// TODO: It could be function.
			variablePtr := scope.FindVariable(*token.Symbol)
			log.Fatal(variablePtr)
			variable := scope.Block.NewLoad(variablePtr)
			expressionIR.Value = variable
		}
	}
}

// Evaluate evalutes the expressionIR and generates the IR.
func (expressionIR *ExpressionIR) Evaluate(scope *Scope) value.Value {
	if expressionIR.Left == nil && expressionIR.Right == nil {
		return expressionIR.Value
	} else {
		value1 := expressionIR.Left.Evaluate(scope)
		value2 := expressionIR.Right.Evaluate(scope)
		result := OperationEvaluateFunc[expressionIR.OperationType](scope, value1, value2)
		return result
	}
}
