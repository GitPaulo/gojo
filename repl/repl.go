package repl

import (
	"errors"
	"fmt"
	"gojo/interpreter"
	"gojo/lexer"
	"gojo/parser"
	"os"
	"strings"

	"github.com/peterh/liner"
)

func StartREPL(inter *interpreter.Interpreter) {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	line.SetMultiLineMode(true)

	// Load command history if it exists
	if f, err := os.Open(".gojo_repl_history"); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	fmt.Println("Welcome to Gojo REPL")
	fmt.Println("Type 'exit' or 'quit' to end the REPL.")

	for {
		input, err := line.Prompt(">>> ")
		if errors.Is(err, liner.ErrPromptAborted) {
			break
		} else if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		// Trim the input
		input = strings.TrimSpace(input)

		// Check for exit command
		if input == "exit" || input == "quit" {
			break
		}

		// Add input to history
		line.AppendHistory(input)

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

	// Save command history
	if f, err := os.Create(".gojo_repl_history"); err != nil {
		fmt.Println("Error saving history:", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}
