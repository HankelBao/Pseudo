package compiler

import (
	"github.com/alecthomas/participle/lexer"
)

type Ast struct {
	Instructions []*Instruction `(@@)+`
}

type Instruction struct {
	Pos             lexer.Position
	Output          *InstOutput          ` @@`
	PrintfD         *InstPrintfD         `|@@`
	PrintfF         *InstPrintfF         `|@@`
	DeclareVariable *InstDeclareVariable `|@@`
	Assignment      *InstAssignment      `|@@`
	ConditionBr     *InstConditionBr     `|@@`
	NullLine        *string              `|@EOL`
}

/*
Example:
	OUTPUT "Hello World!\n"
*/
type InstOutput struct {
	Pos     lexer.Position
	Content Expression `"OUTPUT" @@ EOL`
}

/*
Declare a variable.

Example:
	DECLARE a : INT
*/
type InstDeclareVariable struct {
	Pos  lexer.Position
	Name string       `"DECLARE" @Ident`
	Type VariableType `":" @@ EOL`
}

/*
Assignment

Example:
	a <- 1
*/
type InstAssignment struct {
	Pos   lexer.Position
	Left  Key        `@@ "<"`
	Right Expression `"-" @@ EOL`
}

/*
Output a expression of INT as final type for debug usage.

Example:
	PrintfS 1
*/
type InstPrintfD struct {
	Pos     lexer.Position
	Content Expression `"PrintfD" @@ EOL`
}

type InstPrintfF struct {
	Pos     lexer.Position
	Content Expression `"PrintfF" @@ EOL`
}

type InstConditionBr struct {
	Pos       lexer.Position
	Condition Expression `"IF" @@ EOL`
	TrueBr    Ast        `"THEN" EOL @@`
	FalseBr   *Ast       `("ELSE" @@)?`
	END       string     `"ENDIF" EOL`
}

type VariableType struct {
	Pos    lexer.Position
	Int    *string `  @"INT"`
	REAL   *string `| @"REAL"`
	CUSTOM *string `| @Ident`
}

type Key struct {
	Pos    lexer.Position
	Tokens []*KeyToken `@@+`
}

type KeyToken struct {
	Pos lexer.Position

	Symbol *string `@Ident`

	Dot          *string `| @"."`
	LeftBracket  *string `| @"["`
	RightBracket *string `| @"]"`
}

type Expression struct {
	//Pos        lexer.Position
	Comparison Comparison `@@`
}

type Comparison struct {
	//Pos   lexer.Position
	Head  Addition        `@@`
	Items []*OpComparison `(@@)*`
}

type OpComparison struct {
	//Pos lexer.Position
	Operator string   `@("<" ">" | "=" | "<" "=" | ">" "=" | "<" | ">")`
	Item     Addition `@@`
}

type Addition struct {
	//Pos   lexer.Position
	Head  Multiplication `@@`
	Items []*OpAddition  `(@@)*`
}

type OpAddition struct {
	Operator string         `@("+"|"-")`
	Item     Multiplication `@@`
}

type Multiplication struct {
	//Pos   lexer.Position
	Head  Unary               `@@`
	Items []*OpMultiplication `(@@)*`
}

type OpMultiplication struct {
	Operator string `@("*"|"/")`
	Item     Unary  `@@`
}

type Unary struct {
	//Pos     lexer.Position
	Not      *Unary   `  "!" @@`
	Opposite *Unary   `| "-" @@`
	Primary  *Primary `| @@`
}

type Primary struct {
	//Pos lexer.Position

	Constant      *Constant   ` @@`
	Symbol        *string     `| @Ident`
	Subexpression *Expression `| "(" @@ ")"`
}

type Constant struct {
	//Pos     lexer.Position
	VString *string  `  @String`
	VReal   *float64 `| @Float`
	VInt    *int64   `| @Int`
}
