package compiler

import (
	"io"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
	"log"
)

func Parse(f io.Reader) *Ast{
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

	parser := participle.MustBuild(&Ast{},
		participle.Lexer(pseLexer),
		participle.Unquote("String"),
		participle.UseLookahead(2),
		participle.Elide("Whitespace"),
	)
	ast := &Ast{}
	parse_err := parser.Parse(f, ast)
	if parse_err != nil {
		log.Fatal(parse_err)
	}
	return ast
}

