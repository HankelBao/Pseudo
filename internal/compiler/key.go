package compiler

import (
	//"github.com/llir/llvm/ir"
	//"github.com/llir/llvm/ir/constant"
	//"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	//"log"
)

// Locate tries to locate the key to a value via the scope given.
func (key *Key) Locate(scope *Scope) value.Value {
	variableName := key.Tokens[0].Symbol
	variable := scope.FindVariable(*variableName)
	return variable
}