package compiler

import (
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// value.Value is a interface,
// so it should not be a ptr here.
// (ir.Func is a struct)
type ScopeVariableMap map[string]value.Value
type ScopeFuncMap map[string]*ir.Func

type Scope struct {
	Module *ir.Module
	Block  *ir.Block

	Variables ScopeVariableMap
	Functions ScopeFuncMap

	// For Pseudocode, Anything in the root level belongs to function main.
	// So variables defined in main block are global variables,
	// variables defined in other blocks are private variables.
	Main bool

	// When the current scope is the same as GlobalScope,
	// it is the root scope.
	// Keep this field for function access and constant definition.
	GlobalScope *Scope
	Parent      *Scope
}

func NewGlobalScope() *Scope {
	scope := Scope{
		Module:      ir.NewModule(),
		Block:       nil,
		Variables:   make(ScopeVariableMap),
		Functions:   make(ScopeFuncMap),
		Main:        false,
		GlobalScope: nil,
		Parent:      nil,
	}
	scope.GlobalScope = &scope
	return &scope
}

func (scope *Scope) NewScope(block *ir.Block) *Scope {
	new_scope := Scope{
		Module:      scope.Module,
		Block:       block,
		Variables:   make(ScopeVariableMap),
		Functions:   nil,
		Main:        false,
		GlobalScope: scope.GlobalScope,
		Parent:      scope,
	}
	return &new_scope
}

func (scope *Scope) RegisterVariable(name string, val value.Value) {
	_, ok := scope.Variables[name]
	if ok {
		log.Fatal("Define Variable more than once in a scope: ", name)
	}
	scope.Variables[name] = val
}

func (scope *Scope) FindVariable(name string) value.Value {
	currentScope := scope
	for {
		val, ok := currentScope.Variables[name]
		if ok {
			return val
		}
		if currentScope.Parent == nil {
			return nil
		}
		currentScope = currentScope.Parent
	}
}

func (scope *Scope) RegisterFunction(name string, value *ir.Func) {
	_, ok := scope.Functions[name]
	if ok {
		log.Fatal("Define Function more than once in a scope: ", name)
	}
	scope.Functions[name] = value
}

func (scope *Scope) FindFunction(name string) *ir.Func {
	// Functions only restored in globalscope.
	currentScope := scope.GlobalScope
	val, ok := currentScope.Functions[name]
	if ok {
		return val
	}
	return nil
}

func (scope *Scope) IsGlobal() bool {
	if scope == scope.GlobalScope {
		return true
	}
	return false
}
