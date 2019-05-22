package compiler

import (
	"log"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	//"github.com/alecthomas/repr"

	"strconv"
)

// Evaluate evalues the expression and generates IR.
func (e *Expression) Evaluate(scope *Scope) value.Value {
	return e.Comparison.Evaluate(scope)
}

// Evaluate evalutes comparison
func (c *Comparison) Evaluate(scope *Scope) value.Value {
	lhsValue := c.Head.Evaluate(scope)
	for _, item := range c.Items {
		lhsValue = item.Evaluate(scope, lhsValue)
	}
	return lhsValue
}

// Evaluate evalutes op based on lhs
func (o *OpComparison) Evaluate(scope *Scope, lhsValue value.Value) value.Value {
	rhsValue := o.Item.Evaluate(scope)
	switch o.Operator {
	case "=":
		return CmpEQEval(scope, lhsValue, rhsValue)
	case "<>":
		return CmpNEEval(scope, lhsValue, rhsValue)
	case ">":
		return CmpGTEval(scope, lhsValue, rhsValue)
	case ">=":
		return CmpGEEval(scope, lhsValue, rhsValue)
	case "<":
		return CmpLTEval(scope, lhsValue, rhsValue)
	case "<=":
		return CmpLEEval(scope, lhsValue, rhsValue)
	}
	log.Fatal("unreachable")
	return nil
}

// Evaluate evalutes addition
func (a *Addition) Evaluate(scope *Scope) value.Value {
	lhsValue := a.Head.Evaluate(scope)
	for _, item := range a.Items {
		lhsValue = item.Evaluate(scope, lhsValue)
	}
	return lhsValue
}

// Evaluate evalutes op based on lhs
func (o *OpAddition) Evaluate(scope *Scope, lhsValue value.Value) value.Value {
	rhsValue := o.Item.Evaluate(scope)
	switch o.Operator {
	case "+":
		return AddEval(scope, lhsValue, rhsValue)
	case "-":
		return MinusEval(scope, lhsValue, rhsValue)
	}
	log.Fatal("unreachable")
	return nil
}

// Evaluate evaluates multiplication
func (m *Multiplication) Evaluate(scope *Scope) value.Value {
	lhsValue := m.Head.Evaluate(scope)
	for _, item := range m.Items {
		lhsValue = item.Evaluate(scope, lhsValue)
	}
	return lhsValue
}

// Evaluate evalutes op based on lhs
func (o *OpMultiplication) Evaluate(scope *Scope, lhsValue value.Value) value.Value {
	rhsValue := o.Item.Evaluate(scope)
	switch o.Operator {
	case "*":
		return MultipleEval(scope, lhsValue, rhsValue)
	case "/":
		return DivideEval(scope, lhsValue, rhsValue)
	}
	return nil
}

// Evaluate evaluates unary
func (u *Unary) Evaluate(scope *Scope) value.Value {
	switch {
	// Not is implemented yet
	case u.Opposite != nil:
		return OppositeEval(scope, u.Opposite.Evaluate(scope))
	case u.Primary != nil:
		return u.Primary.Evaluate(scope)
	}
	log.Fatal("unreachable")
	return nil
}

// Evaluate evaluates primary
func (p *Primary) Evaluate(scope *Scope) value.Value {
	switch {
	case p.Constant != nil:
		return p.Constant.Evaluate(scope)
	case p.Symbol != nil:
		variablePtr := scope.FindVariable(*p.Symbol)
		variable := scope.Block.NewLoad(variablePtr)
		return variable
	case p.Subexpression != nil:
		return p.Subexpression.Evaluate(scope)
	}
	log.Fatal("unreacheable")
	return nil
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
	case c.VBool != nil:
		switch *c.VBool {
		case "TRUE":
			return constant.NewBool(true)
		case "FALSE":
			return constant.NewBool(false)
		}
	default:
		// log.Fatal("Parser: cannot parse value at", c.Pos)
	}
	return nil
}
