package main

import (
	"fmt"
	"gojo/interpreter"
	"gojo/lexer"
	"gojo/parser"
)

func main() {
	input := `
		var x = 5;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
		return
	}

	inter := interpreter.New()
	inter.Interpret(program)
}
