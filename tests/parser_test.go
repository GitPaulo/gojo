package tests

import (
	"fmt"
	"gojo/lexer"
	. "gojo/parser"
	"os"
	"testing"
)

type ParserTestCase struct {
	Name     string
	Expected string
}

func TestParser(t *testing.T) {
	for _, test := range parserTestCases {
		t.Run(
			test.Name,
			func(t *testing.T) {
				CompareParserOutput(t, test)
			},
		)
	}
}

func CompareParserOutput(t *testing.T, test ParserTestCase) {
	const testDataDir = "data/parser"
	filePath := fmt.Sprintf("%s/%s.js", testDataDir, test.Name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Could not read file: %q", filePath)
	}

	lex := lexer.New(string(data))
	parser := New(lex)
	program := parser.ParseProgram()
	if program.String() != test.Expected {
		t.Fatalf("\nExpected: %v\nReceived: %v\n", test.Expected, program.String())
	}
}

var parserTestCases = []ParserTestCase{
	{
		Name:     "Test1",
		Expected: `Program(VariableDeclaration(var Identifier(x) = IntegerLiteral(5))VariableDeclaration(var Identifier(y) = IntegerLiteral(10))VariableDeclaration(var Identifier(z) = BinaryExpression(Identifier(x) + Identifier(y))))`,
	},
	{
		Name:     "Test2",
		Expected: `Program(VariableDeclaration(var Identifier(a) = BooleanLiteral(true))VariableDeclaration(var Identifier(b) = BooleanLiteral(false))VariableDeclaration(var Identifier(c) = BinaryExpression(Identifier(a) && Identifier(b))))`,
	},
}
