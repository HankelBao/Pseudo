package main

import (
	"fmt"
	"github.com/HankelBao/GoPse/internal/compiler"
	"github.com/alecthomas/repr"
	"log"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Filename required")
	}
	filename := os.Args[1]
	in, err := os.Open(filename)
	if err != nil {
		log.Fatal("Could not open file")
	}

	ast := compiler.Parse(in)
	repr.Println(ast)

	m := compiler.Compile(ast)
	fmt.Println(m)

	os.MkdirAll("./tmp", os.ModePerm)
	out, _ := os.Create("./tmp/test.ll")
	out.Write([]byte(m.String()))
	out.Close()

	cmd_output, cmd_err := exec.Command("clang", "./tmp/test.ll", "-o", "./tmp/test.exe").Output()
	if cmd_err != nil {
		log.Fatal(cmd_err)
	}
	log.Println(cmd_output)
}
