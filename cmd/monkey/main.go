package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xPoppa/interpreter/evaluator"
	"github.com/xPoppa/interpreter/lexer"
	"github.com/xPoppa/interpreter/object"
	"github.com/xPoppa/interpreter/parser"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("Need the path of the source file\n")
		os.Exit(1)
	}
	readAndEvaluateFile(args[1])
}

func readAndEvaluateFile(path string) {
	f, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Cannot find file: %s", path)
	}
	if err != nil {
		log.Fatalf("Couldn't read file with err: %s", err)
	}
	file := string(f)
	env := object.NewEnvironment()
	l := lexer.New(file)
	p := parser.New(l)

	program := p.ParseProgram()
	evaluated := evaluator.Eval(program, env)
	fmt.Println(evaluated.Inspect())
}
