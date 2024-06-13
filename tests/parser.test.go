package tests

import (
	"fmt"
	"gojo/lexer"
	"gojo/parser"
)

func RunParserTests() {
	fmt.Println("=== Parser Tests ===")
	tests := []struct {
		input    string
		expected string
	}{
		{
			input: `
			var x = 5;
			var y = 10;
			var z = x + y;
			`,
			expected: `Program(VariableDeclaration(var x = 5)VariableDeclaration(var y = 10)VariableDeclaration(var z = BinaryExpression(x + y)))`,
		},
		{
			input: `
			var a = true;
			var b = false;
			var c = a && b;
			`,
			expected: `Program(VariableDeclaration(var a = true)VariableDeclaration(var b = false)VariableDeclaration(var c = BinaryExpression(a && b))`,
		},
	}

	for i, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		if program.String() != test.expected {
			fmt.Printf("❌ Test %d failed. Expected: %s, Got: %s\n", i+1, test.expected, program.String())
		} else {
			fmt.Printf("✅ Test %d passed.\n", i+1)
		}
	}
}
