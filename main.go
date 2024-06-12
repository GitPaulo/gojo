package main

import (
	"fmt"
	"gojo/config"
	"gojo/interpreter"
	"gojo/lexer"
	"gojo/parser"
	"gojo/tests"
	"os"
)

func main() {
	env := config.LoadConfig()
	if env.Test {
		fmt.Println("=== Lexer Tests ===")
		tests.RunLexerTests()
		fmt.Println("=== Parser Tests ===")
		tests.RunParserTests()
		fmt.Println("=== Interpreter Tests ===")
		tests.RunInterpreterTests()
		fmt.Println("Tests ran. Exiting...")
		return
	}

	// Input through input_program.js
	inputFile := env.InputFile
	if inputFile == "" {
		inputFile = "input_program.js"
	}
	var input string
	var err error
	input, err = readFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}
	printInput(input)

	// Initialize lexer, parser, and interpreter
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

func printInput(input string) {
	fmt.Println("Input:")
	fmt.Println("------")
	fmt.Println(input)
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

func readFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
