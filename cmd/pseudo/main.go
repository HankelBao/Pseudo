package main

import (
	"fmt"
	"github.com/HankelBao/Pseudo/internal/compiler"
	"github.com/alecthomas/repr"
	"log"
	"os"
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

	/*cmdOutput, cmdErr := exec.Command("clang", "./tmp/test.ll", "-o", "./tmp/test.exe").Output()
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
	log.Println(cmdOutput)*/
}
