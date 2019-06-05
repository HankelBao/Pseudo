package compiler

import (
	"github.com/alecthomas/participle/lexer"
)

// Ast is the instructions assembly
type Ast struct {
	Instructions []*Instruction `(@@)+`
}

// Instruction matches all kinds of instructions
type Instruction struct {
	Pos             lexer.Position
	Output          *InstOutput          ` @@`
	Input           *InstInput           `|@@`
	Call            *InstCall            `|@@`
	PrintfD         *InstPrintfD         `|@@`
	PrintfF         *InstPrintfF         `|@@`
	DeclareVariable *InstDeclareVariable `|@@`
	ConditionBr     *InstConditionBr     `|@@`
	While           *InstWhile           `|@@`
	Repeat          *InstRepeat          `|@@`
	Assignment      *InstAssignment      `|@@`
	NullLine        *string              `|@EOL`
}

// InstOutput outputs string
// Example:
// 	OUTPUT "Hello World!\n"
type InstOutput struct {
	Pos     lexer.Position
	Content Expression `"OUTPUT" @@ EOL`
}

type InstInput struct {
	Pos     lexer.Position
	Content Key `"INPUT" @@ EOL`
}

// InstDeclareVariable declares a variable.
//
// Example:
// 	DECLARE a : INT
type InstDeclareVariable struct {
	Pos  lexer.Position
	Name string       `"DECLARE" @Ident`
	Type VariableType `":" @@ EOL`
}

// InstAssignment assigns a variable the value of an expression
//
// Example:
// 	a <- 1
type InstAssignment struct {
	Pos   lexer.Position
	Left  Key        `@@ "<"`
	Right Expression `"-" @@ EOL`
}

// InstPrintfD outputs a expression of INT for debug usage
//
// Example:
// 	PrintfS 1
type InstPrintfD struct {
	Pos     lexer.Position
	Content Expression `"PrintfD" @@ EOL`
}

// InstPrintfF outputs a expression of REAL for debug usage
//
// Example:
// 	PrintF 1.0
type InstPrintfF struct {
	Pos     lexer.Position
	Content Expression `"PrintfF" @@ EOL`
}

// InstConditionBr creates if..then..else...
//
// Example:
// 	IF 1==1
// 	  THEN
//	    OUTPUT "TURE"
// 	ENDIF
type InstConditionBr struct {
	Pos       lexer.Position
	Condition Expression `"IF" @@ EOL`
	TrueBr    Ast        `"THEN" EOL @@`
	FalseBr   *Ast       `("ELSE" EOL @@)?`
	END       string     `"ENDIF" EOL`
}

// InstWhile creates while block
//
// Example:
//	WHILE 1=1 DO
// 		OUTPUT "Hi"
// 	ENDWHILE
type InstWhile struct {
	Pos       lexer.Position
	Condition Expression `"WHILE" @@ "DO" EOL`
	Body      Ast        `@@`
	END       string     `"ENDWHILE" EOL`
}

// InstWhile creates repeat block
//
// Example:
//	REPEAT
// 		OUTPUT "HI"
//	UNTIL 1=1
type InstRepeat struct {
	Pos       lexer.Position
	Head      string     `"REPEAT" EOL`
	Body      Ast        `@@`
	Condition Expression `"UNTIL" @@ EOL`
}

// InstCall creates call block
//
// Example:
// 	CALL puts("Hello World!")
type InstCall struct {
	Pos      lexer.Position
	Function *FunctionCall `"CALL" @@ EOL`
}

// VariableType matches the variable type of declaration
type VariableType struct {
	Pos    lexer.Position
	BOOL   *string `  @"BOOL"`
	Int    *string `| @"INT"`
	REAL   *string `| @"REAL"`
	CUSTOM *string `| @Ident`
}

// Key is an assignable terminal
type Key struct {
	Pos       lexer.Position
	Variables []*Variable `@@ ("." @@)*`
}

type Variable struct {
	Pos        lexer.Position
	Name       string `@Ident`
	ArrayIndex *int   `("[" @Int "]")?`
}

// KeyToken is the lexers of Key
// TODO: change into expression-like handling
type KeyToken struct {
	Pos lexer.Position

	Symbol *string `@Ident`

	Dot          *string `| @"."`
	LeftBracket  *string `| @"["`
	RightBracket *string `| @"]"`
}

// Expression is expression of an value
type Expression struct {
	//Pos        lexer.Position
	Comparison Comparison `@@`
}

// Comparison compares two or more values
type Comparison struct {
	//Pos   lexer.Position
	Head  Addition        `@@`
	Items []*OpComparison `(@@)*`
}

// OpComparison makes a comparison with another value
type OpComparison struct {
	//Pos lexer.Position
	Operator string   `@("<" ">" | "=" | "<" "=" | ">" "=" | "<" | ">")`
	Item     Addition `@@`
}

// Addition adds two or more values
type Addition struct {
	//Pos   lexer.Position
	Head  Multiplication `@@`
	Items []*OpAddition  `(@@)*`
}

// OpAddition adds another value
type OpAddition struct {
	Operator string         `@("+"|"-")`
	Item     Multiplication `@@`
}

// Multiplication multiples two values
type Multiplication struct {
	//Pos   lexer.Position
	Head  Unary               `@@`
	Items []*OpMultiplication `(@@)*`
}

// OpMultiplication multiples with another value
type OpMultiplication struct {
	Operator string `@("*"|"/")`
	Item     Unary  `@@`
}

// Unary gives not or opposite
type Unary struct {
	//Pos     lexer.Position
	Not      *Unary   `  "!" @@`
	Opposite *Unary   `| "-" @@`
	Primary  *Primary `| @@`
}

// Primary is the smallest universal unit in an expression
type Primary struct {
	//Pos lexer.Position

	Constant      *Constant     `  @@`
	Function      *FunctionCall `| @@`
	Key           *Key          `| @@`
	Subexpression *Expression   `| "(" @@ ")"`
}

// Constant shows direct value
type Constant struct {
	//Pos     lexer.Position
	VBool   *string  `  @("TRUE"|"FALSE")`
	VString *string  `| @String`
	VReal   *float64 `| @Float`
	VInt    *int64   `| @Int`
}

type FunctionCall struct {
	Name   string        `@Ident`
	Params []*Expression `"(" (@@ ("," @@)*)? ")"`
}
