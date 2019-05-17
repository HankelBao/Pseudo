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
	Pos    lexer.Position
	Tokens []*ExpressionToken `@@+`
}

// In order to simplify the tokens and improve the error report,
// Expression Tokens are just lexers.
//
// (...) could represent prioritized expression or function params
// [...] could only represent index (expression)
type ExpressionToken struct {
	Pos lexer.Position

	BasicOp     *string `  @("+" | "-" | "*" | "/")`
	CmpOp       *string `| @("<" ">" | "=" | "<" "=" | ">" "=" | "<" | ">")`
	Parenthesis *string `| @("(" | ")")`
	Bracket     *string `| @("[" | "]")`

	Constant *Constant `| @@`
	Symbol   *string   `| @Ident`

	ParamDivier *string `| @","`
}

type Constant struct {
	Pos     lexer.Position
	VString *string  `  @String`
	VReal   *float64 `| @Float`
	VInt    *int64   `| @Int`
}
