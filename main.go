package main

import (
	"fmt"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/repr"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"os"
)

type Entry struct {
	Pos          lexer.Position
	Instructions []*Instruction `(@@)+`
}

type Instruction struct {
	Pos    lexer.Position
	Output *InstOutput `@@`
}

type InstOutput struct {
	Pos     lexer.Position
	Title   string     `@"OUTPUT"`
	Content Expression `@@`
}

type Expression struct {
	Pos    lexer.Position
	Tokens []*ExpressionToken `(@@)+`
}

type ExpressionToken struct {
	Pos          lexer.Position
	Add          *string `  @"+"`
	Minus        *string `| @"-"`
	Times        *string `| @"*"`
	Divide       *string `| @"/"`
	Value        *Value  `| @@`
	Variable     *string `| @Ident`
	LeftBracket  *string `| @"("`
	RightBracket *string `| @")"`
}

type Value struct {
	Pos     lexer.Position
	VString *string  `  @String`
	VReal   *float64 `| @Float`
	VInt    *int     `| @Int`
}

func main() {
	parser, _ := participle.Build(&Entry{})
	ast := &Entry{}
	in, _ := os.Open("./test.pse")
	parser.Parse(in, ast)
	in.Close()

	repr.Println(ast)

	m := ir.NewModule()

	puts := m.NewFunc("puts", &types.IntType{BitSize: 32}, ir.NewParam("", &types.PointerType{ElemType: types.I8}))
	helloworld_str := constant.NewCharArray([]byte("Hello World"))
	zero := constant.NewInt(types.I32, 0)
	helloworld_def := m.NewGlobalDef("str", helloworld_str)

	main := m.NewFunc("main", &types.IntType{BitSize: 32})
	entry := main.NewBlock("")
	tmp := entry.NewBitCast(helloworld_def, &types.PointerType{ElemType: &types.IntType{BitSize: 8}})
	entry.NewCall(puts, tmp)
	entry.NewRet(zero)

	fmt.Println(m)
	out, _ := os.Create("./test.ll")
	out.Write([]byte(m.String()))
	out.Close()

}
