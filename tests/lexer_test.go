package tests

import (
	"fmt"
	. "gojo/lexer"
	"os"
	"testing"
)

type LexerTestCase struct {
	Name     string
	Expected []GojoToken
}

func TestLexer(t *testing.T) {
	for _, test := range lexerTestCases {
		t.Run(
			test.Name,
			func(t *testing.T) {
				CompareLexerOutput(t, test)
			},
		)
	}
}

func CompareLexerOutput(t *testing.T, test LexerTestCase) {
	const testDataDir = "data/lexer"
	filePath := fmt.Sprintf("%s/%s.js", testDataDir, test.Name)
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Could not read file: %q", filePath)
	}

	lexer := New(string(data))
	for i, expectedToken := range test.Expected {
		token := lexer.NextToken()
		if token.Type != expectedToken.Type || token.Text != expectedToken.Text {
			t.Fatalf("Token %2d: \nExpected: %v\nReceived: %v\n", i, expectedToken, token)
		} else {
			t.Logf("Token %2d: %v", i, token)
		}
	}
}

var lexerTestCases = []LexerTestCase{
	// Test variable declarations and basic arithmetic operations
	{
		Name: "Test1",
		Expected: []GojoToken{
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "x"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenLiterals["number"], Text: "5"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "y"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenLiterals["number"], Text: "10"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "z"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenText["identifier"], Text: "x"},
			{Type: TokenOperators["+"], Text: "+"},
			{Type: TokenText["identifier"], Text: "y"},
			{Type: TokenPunctuation[";"], Text: ";"},
		},
	},
	// Test function declaration and call
	{
		Name: "Test2",
		Expected: []GojoToken{
			{Type: TokenKeywords["function"], Text: "function"},
			{Type: TokenText["identifier"], Text: "add"},
			{Type: TokenPunctuation["("], Text: "("},
			{Type: TokenText["identifier"], Text: "a"},
			{Type: TokenPunctuation[","], Text: ","},
			{Type: TokenText["identifier"], Text: "b"},
			{Type: TokenPunctuation[")"], Text: ")"},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenKeywords["return"], Text: "return"},
			{Type: TokenText["identifier"], Text: "a"},
			{Type: TokenOperators["+"], Text: "+"},
			{Type: TokenText["identifier"], Text: "b"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenPunctuation["}"], Text: "}"},
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "result"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenText["identifier"], Text: "add"},
			{Type: TokenPunctuation["("], Text: "("},
			{Type: TokenLiterals["number"], Text: "5"},
			{Type: TokenPunctuation[","], Text: ","},
			{Type: TokenLiterals["number"], Text: "10"},
			{Type: TokenPunctuation[")"], Text: ")"},
			{Type: TokenPunctuation[";"], Text: ";"},
		},
	},
	// Test for loop and conditionals
	{
		Name: "Test3",
		Expected: []GojoToken{
			{Type: TokenKeywords["for"], Text: "for"},
			{Type: TokenPunctuation["("], Text: "("},
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "i"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenLiterals["number"], Text: "0"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenText["identifier"], Text: "i"},
			{Type: TokenOperators["<"], Text: "<"},
			{Type: TokenLiterals["number"], Text: "10"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenText["identifier"], Text: "i"},
			{Type: TokenOperators["++"], Text: "++"},
			{Type: TokenPunctuation[")"], Text: ")"},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenKeywords["if"], Text: "if"},
			{Type: TokenPunctuation["("], Text: "("},
			{Type: TokenText["identifier"], Text: "i"},
			{Type: TokenOperators["%"], Text: "%"},
			{Type: TokenLiterals["number"], Text: "2"},
			{Type: TokenOperators["=="], Text: "=="},
			{Type: TokenLiterals["number"], Text: "0"},
			{Type: TokenPunctuation[")"], Text: ")"},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenKeywords["continue"], Text: "continue"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenPunctuation["}"], Text: "}"},
			{Type: TokenKeywords["else"], Text: "else"},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenKeywords["break"], Text: "break"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenPunctuation["}"], Text: "}"},
			{Type: TokenPunctuation["}"], Text: "}"},
		},
	},
	// Test string literals and escape sequences
	{
		Name: "Test4",
		Expected: []GojoToken{
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "str"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenLiterals["string"], Text: "Hello, world!\n"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "escapedStr"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenLiterals["string"], Text: "This is a \"quoted\" string."},
			{Type: TokenPunctuation[";"], Text: ";"},
		},
	},
	// Test boolean literals and logical operators
	{
		Name: "Test5",
		Expected: []GojoToken{
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "a"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenKeywords["true"], Text: "true"},
			{Type: TokenOperators["&&"], Text: "&&"},
			{Type: TokenKeywords["false"], Text: "false"},
			{Type: TokenOperators["||"], Text: "||"},
			{Type: TokenOperators["!"], Text: "!"},
			{Type: TokenKeywords["true"], Text: "true"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "b"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenOperators["!"], Text: "!"},
			{Type: TokenOperators["!"], Text: "!"},
			{Type: TokenText["identifier"], Text: "a"},
			{Type: TokenPunctuation[";"], Text: ";"},
		},
	},
	// Test object literals
	{
		Name: "Test6",
		Expected: []GojoToken{
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "obj"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenText["identifier"], Text: "key1"},
			{Type: TokenPunctuation[":"], Text: ":"},
			{Type: TokenLiterals["string"], Text: "value1"},
			{Type: TokenPunctuation[","], Text: ","},
			{Type: TokenText["identifier"], Text: "key2"},
			{Type: TokenPunctuation[":"], Text: ":"},
			{Type: TokenLiterals["number"], Text: "42"},
			{Type: TokenPunctuation[","], Text: ","},
			{Type: TokenText["identifier"], Text: "key3"},
			{Type: TokenPunctuation[":"], Text: ":"},
			{Type: TokenKeywords["true"], Text: "true"},
			{Type: TokenPunctuation["}"], Text: "}"},
			{Type: TokenPunctuation[";"], Text: ";"},
		},
	},
	// Test array literals and access
	{
		Name: "Test7",
		Expected: []GojoToken{
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "arr"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenPunctuation["["], Text: "["},
			{Type: TokenLiterals["number"], Text: "1"},
			{Type: TokenPunctuation[","], Text: ","},
			{Type: TokenLiterals["number"], Text: "2"},
			{Type: TokenPunctuation[","], Text: ","},
			{Type: TokenLiterals["number"], Text: "3"},
			{Type: TokenPunctuation["]"], Text: "]"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "first"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenText["identifier"], Text: "arr"},
			{Type: TokenPunctuation["["], Text: "["},
			{Type: TokenLiterals["number"], Text: "0"},
			{Type: TokenPunctuation["]"], Text: "]"},
			{Type: TokenPunctuation[";"], Text: ";"},
		},
	},
	// Test class declaration
	{
		Name: "Test8",
		Expected: []GojoToken{
			{Type: TokenKeywords["class"], Text: "class"},
			{Type: TokenText["identifier"], Text: "Person"},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenText["identifier"], Text: "constructor"},
			{Type: TokenPunctuation["("], Text: "("},
			{Type: TokenText["identifier"], Text: "name"},
			{Type: TokenPunctuation[")"], Text: ")"},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenKeywords["this"], Text: "this"},
			{Type: TokenPunctuation["."], Text: "."},
			{Type: TokenText["identifier"], Text: "name"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenText["identifier"], Text: "name"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenPunctuation["}"], Text: "}"},
			{Type: TokenText["identifier"], Text: "getName"},
			{Type: TokenPunctuation["("], Text: "("},
			{Type: TokenPunctuation[")"], Text: ")"},
			{Type: TokenPunctuation["{"], Text: "{"},
			{Type: TokenKeywords["return"], Text: "return"},
			{Type: TokenKeywords["this"], Text: "this"},
			{Type: TokenPunctuation["."], Text: "."},
			{Type: TokenText["identifier"], Text: "name"},
			{Type: TokenPunctuation[";"], Text: ";"},
			{Type: TokenPunctuation["}"], Text: "}"},
			{Type: TokenPunctuation["}"], Text: "}"},
		},
	},
	// Test regular expressions
	{
		Name: "Test9",
		Expected: []GojoToken{
			{Type: TokenKeywords["var"], Text: "var"},
			{Type: TokenText["identifier"], Text: "regex"},
			{Type: TokenOperators["="], Text: "="},
			{Type: TokenLiterals["regexp"], Text: "/ab+c/"},
			{Type: TokenPunctuation[";"], Text: ";"},
		},
	},
}
