package compiler

import (
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)


// ScopeVariableMap is a map to store all the variables in a scope.
// value.Value is an interface,
// so it should not be a ptr here.
type ScopeVariableMap map[string]value.Value
// ScopeFuncMap is a map to store all the functions in a scope.
// ir.Func is a struct
type ScopeFuncMap map[string]*ir.Func

// Scope keep track of all the informations in a block/sub-block.
type Scope struct {
	Module *ir.Module
	Func   *ir.Func
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

// NewGlobalScope creates a global scope.
// There should be only one global scope.
func NewGlobalScope() *Scope {
	scope := Scope{
		Module:      ir.NewModule(),
		Func:        nil,
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

// NewFuncScope creates a function scope under global scope.
func (scope *Scope) NewFuncScope(function *ir.Func) *Scope {
	newScope := Scope{
		Module:      scope.Module,
		Func:        function,
		Block:       function.NewBlock(""),
		Variables:   make(ScopeVariableMap),
		Functions:   nil,
		Main:        false,
		GlobalScope: scope.GlobalScope,
		Parent:      scope,
	}
	return &newScope
}

// NewScope creates a new scope under the given scope.
func (scope *Scope) NewScope(block *ir.Block) *Scope {
	if scope.Func == nil {
		log.Fatal("Cannot new a scope here!")
	}
	newScope := Scope{
		Module:      scope.Module,
		Func:        scope.Func,
		Block:       block,
		Variables:   make(ScopeVariableMap),
		Functions:   nil,
		Main:        false,
		GlobalScope: scope.GlobalScope,
		Parent:      scope,
	}
	return &newScope
}

// RegisterVariable register a variable to the current scope for further usages.
func (scope *Scope) RegisterVariable(name string, val value.Value) {
	_, ok := scope.Variables[name]
	if ok {
		log.Fatal("Define Variable more than once in a scope: ", name)
	}
	scope.Variables[name] = val
}

// FindVariable locates the variable registered.
// If the variable is not found, nil would be returned.
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

// RegisterFunction registers a function to the current scope for further usages.
// Functions should be registered to global scope only!
func (scope *Scope) RegisterFunction(name string, value *ir.Func) {
	_, ok := scope.Functions[name]
	if ok {
		log.Fatal("Define Function more than once in a scope: ", name)
	}
	scope.Functions[name] = value
}

// FindFunction locates the function in the current scope.
func (scope *Scope) FindFunction(name string) *ir.Func {
	// Functions only restored in globalscope.
	currentScope := scope.GlobalScope
	val, ok := currentScope.Functions[name]
	if ok {
		return val
	}
	return nil
}

// IsGlobal checks if the current scope is the global scope.
func (scope *Scope) IsGlobal() bool {
	if scope == scope.GlobalScope {
		return true
	}
	return false
}
