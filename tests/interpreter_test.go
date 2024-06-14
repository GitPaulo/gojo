package tests

import (
	"fmt"
	. "gojo/interpreter"
	"gojo/lexer"
	"gojo/parser"
	"os"
	"testing"
)

type InterpreterTestCase struct {
	Name     string
	Expected map[string]interface{}
}

func TestInterpreter(t *testing.T) {
	for _, test := range interpreterTestCases {
		t.Run(
			test.Name,
			func(t *testing.T) {
				CompareInterpreterOutput(t, test)
			},
		)
	}
}

func CompareInterpreterOutput(t *testing.T, test InterpreterTestCase) {
	const testDataDir = "data/interpreter"
	filePath := fmt.Sprintf("%s/%s.js", testDataDir, test.Name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Could not read file: %q", filePath)
	}

	l := lexer.New(string(data))
	p := parser.New(l)
	program := p.ParseProgram()
	interpreter := New()
	interpreter.Interpret(program)

	for key, expectedValue := range test.Expected {
		if value, exists := interpreter.Env[key]; !exists || value != expectedValue {
			t.Errorf("Variable: %2s\nExpected: %2v\nReceived: %2v", key, expectedValue, value)
		} else {
			t.Logf("Variable: %2s = %2v", key, value)
		}
	}
}

var interpreterTestCases = []InterpreterTestCase{
	{
		Name: "Test1",
		Expected: map[string]interface{}{
			"x": int64(5),
			"y": int64(10),
			"z": int64(15),
		},
	},
	{
		Name: "Test2",
		Expected: map[string]interface{}{
			"a": true,
			"b": false,
			"c": false,
		},
	},
}
