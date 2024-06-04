package main

import (
	"fmt"
	"gojo/config"
	"gojo/interpreter"
	"gojo/lexer"
	"gojo/parser"
)

func main() {
	// Define the input program
	input := `
		var x = 5;
		var y = 10;
		var z = x + y;
	`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if config.LoadConfig().Verbose {
		printProgramDetails(program)
	}

	// Check for parser errors
	if len(p.Errors()) != 0 {
		printParserErrors(p)
		return
	}

	inter := interpreter.New()
	inter.Interpret(program)
}

func printProgramDetails(program *parser.Program) {
	fmt.Println("Program:")
	fmt.Printf("  Statements: (%d elements)\n", len(program.Statements))
	for _, stmt := range program.Statements {
		fmt.Printf("    - %s\n", stmt)
	}
	fmt.Printf("  Start: %d\n", program.Start)
	fmt.Printf("  End: %d\n", program.End)
}

func printParserErrors(p *parser.Parser) {
	fmt.Println("Parser Errors:")
	for _, err := range p.Errors() {
		fmt.Println(err)
	}
}
