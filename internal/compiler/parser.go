package compiler

import (
	"io"
	"github.com/alecthomas/participle"
	"log"
)

func Parse(f io.Reader) *Ast{
	parser, err := participle.Build(&Ast{})
	if err != nil {
		log.Fatalln("Could not create parser")
	}
	ast := &Ast{}
	parser.Parse(f, ast)
	return ast
}

