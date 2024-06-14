package repl

import (
	"bufio"
	"fmt"
	"gojo/interpreter"
	"gojo/lexer"
	"gojo/parser"
	"os"
)

func StartREPL(inter *interpreter.Interpreter) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to Gojo REPL")
	fmt.Println("Type 'exit' or 'quit' to end the REPL.")
	for {
		fmt.Print(">>> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Trim the input
		input = input[:len(input)-1]

		// Check for exit command
		if input == "exit" || input == "quit" {
			break
		}

		// Create a lexer and parser for the input
		lex := lexer.New(input)
		p := parser.New(lex)

		// Parse the input to create a program AST
		program := p.ParseProgram()

		// Check for parsing errors
		if len(p.Errors()) > 0 {
			for _, e := range p.Errors() {
				fmt.Println("Parsing error:", e)
			}
			continue
		}

		// Interpret the parsed program
		inter.InterpretREPL(program)
	}
}
