package compiler

import (
	"io"
	"log"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

// Parse parses the input into ast.
func Parse(f io.Reader) *Ast {
	pseLexer := lexer.Must(ebnf.New(`
		Ident = (alpha | "_") { "_" | alpha | digit } .
		String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
		Float = digit {digit} "." digit {digit} .
		Int = digit {digit} .
		EOL = [ "\r" ] "\n" .
		Whitespace = ( " " | "\t" ) { " " | "\t" } .
		Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .

		alpha = "a"…"z" | "A"…"Z" .
		digit = "0"…"9" .
		any = "\u0000"…"\uffff" .
	`))

	parser, parserErr := participle.Build(&Ast{},
		participle.Lexer(pseLexer),
		participle.Unquote("String"),
		//participle.UseLookahead(0),
		participle.Elide("Whitespace"),
	)
	if parserErr != nil {
		log.Fatal(parserErr)
	}

	ast := &Ast{}
	parseErr := parser.Parse(f, ast)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	return ast
}
