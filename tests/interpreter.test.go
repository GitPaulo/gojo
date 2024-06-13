package tests

import (
	"fmt"
	"gojo/interpreter"
	"gojo/lexer"
	"gojo/parser"
)

func RunInterpreterTests() {
	fmt.Println("=== Interpreter Tests ===")
	tests := []struct {
		input    string
		expected map[string]interface{}
	}{
		{
			input: `
			var x = 5;
			var y = 10;
			var z = x + y;
			`,
			expected: map[string]interface{}{
				"x": int64(5),
				"y": int64(10),
				"z": int64(15),
			},
		},
		{
			input: `
			var a = true;
			var b = false;
			var c = a && b;
			`,
			expected: map[string]interface{}{
				"a": true,
				"b": false,
				"c": false,
			},
		},
	}

	for i, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		interp := interpreter.New()
		interp.Interpret(program)

		testPassed := true
		for key, expectedValue := range test.expected {
			if value, exists := interp.Env[key]; !exists || value != expectedValue {
				fmt.Printf("❌ Test %d failed. For variable '%s', expected: %v, got: %v\n", i+1, key, expectedValue, value)
				testPassed = false
			}
		}
		if testPassed {
			fmt.Printf("✅ Test %d passed.\n", i+1)
		}
	}
}
