package compiler

import (
	"github.com/alecthomas/participle/lexer"
)

const EOLTag = `("\n"|"\r") ("\n"|"\r")*`

type Ast struct {
	Instructions []*Instruction `(@@)+`
}

type Instruction struct {
	Pos             lexer.Position
	Output          *InstOutput          ` @@`
	PrintfS			*InstPrintfS		 `|@@`
	DeclareVariable *InstDeclareVariable `|@@`
	Assignment      *InstAssignment      `|@@`
	NullLine		*string				 `|@EOL`
}

/*
Example:
	OUTPUT "Hello World!\n"
*/
type InstOutput struct {
	Pos     lexer.Position
	Inst    string     `"OUTPUT"`
	Content Expression `@@`
	EOL string `@EOL`
}

/*
Declare a variable.

Example:
	DECLARE a : INT
*/
type InstDeclareVariable struct {
	Pos  lexer.Position
	Name string       `"DECLARE" @Ident`
	Type VariableType `":" @@`
	EOL string `@EOL`
}

/*
Assignment

Example:
	a <- 1
*/
type InstAssignment struct {
	Pos   lexer.Position
	Left  Key        `@@ "<"`
	Right Expression `"-" @@`
	EOL string `@EOL`
}

/*
Output a expression of INT as final type for debug usage.

Example:
	PrintfS 1
*/
type InstPrintfS struct {
	Pos lexer.Position
	Title string `"PrintfS"`
	Content Expression `@@`
	EOL string `@EOL`
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
	Pos    lexer.Position
	Add    *string `  @"+"`
	Minus  *string `| @"-"`
	Times  *string `| @"*"`
	Divide *string `| @"/"`

	CmpEQ *string `| @"="`
	CmpNE *string `| @"!="`
	CmpGR *string `| @">"`
	CmpGE *string `| @">="`
	CmpLT *string `| @"<"`
	CmpLE *string `| @"<="`

	Constant *Constant `| @@`
	Symbol   *string   `| @Ident`

	LeftParenthesis  *string `| @"("`
	RightParenthesis *string `| @")"`
	LeftBracket      *string `| @"["`
	RightBracket     *string `| @"]"`
	ParamDivier      *string `| @","`
}

type Constant struct {
	Pos     lexer.Position
	VString *string  `  @String`
	VReal   *float64 `| @Float`
	VInt    *int64   `| @Int`
}
