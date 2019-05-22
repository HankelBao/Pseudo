package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/HankelBao/Pseudo/internal/compiler"
	"github.com/alecthomas/repr"
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
	// fmt.Println(m)

	os.MkdirAll("./tmp", os.ModePerm)
	out, _ := os.Create("./tmp/test.ll")
	out.Write([]byte(m.String()))
	out.Close()

	var targetName string
	if runtime.GOOS == "windows" {
		targetName = "./tmp/test.exe"
	} else if runtime.GOOS == "linux" {
		targetName = "./tmp/test"
	}

	clangCmdOutput, clangCmdErr := exec.Command("clang", "./tmp/test.ll", "-o", targetName).Output()
	if clangCmdErr != nil {
		log.Fatal(clangCmdErr)
	}
	if string(clangCmdOutput) != "" {
		log.Println(string(clangCmdOutput))
	}

	execCmdOutput, execCmdErr := exec.Command(targetName).Output()
	if execCmdErr != nil {
		log.Fatal(execCmdErr)
	}
	fmt.Println(string(execCmdOutput))
}
