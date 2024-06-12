package tests

import (
	"fmt"
	"gojo/config"
	"gojo/lexer"
)

type TestCase struct {
	Input    string
	Expected []lexer.GojoToken
}

func RunLexerTests() {
	env := config.LoadConfig()
	tests := []TestCase{
		// Test variable declarations and basic arithmetic operations
		{
			Input: `
			var x = 5;
			var y = 10;
			var z = x + y;
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "x"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenLiterals["number"], Text: "5"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "y"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenLiterals["number"], Text: "10"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "z"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenText["identifier"], Text: "x"},
				{Type: lexer.TokenOperators["+"], Text: "+"},
				{Type: lexer.TokenText["identifier"], Text: "y"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
			},
		},
		// Test function declaration and call
		{
			Input: `
			function add(a, b) {
				return a + b;
			}
			var result = add(5, 10);
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["function"], Text: "function"},
				{Type: lexer.TokenText["identifier"], Text: "add"},
				{Type: lexer.TokenPunctuation["("], Text: "("},
				{Type: lexer.TokenText["identifier"], Text: "a"},
				{Type: lexer.TokenPunctuation[","], Text: ","},
				{Type: lexer.TokenText["identifier"], Text: "b"},
				{Type: lexer.TokenPunctuation[")"], Text: ")"},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenKeywords["return"], Text: "return"},
				{Type: lexer.TokenText["identifier"], Text: "a"},
				{Type: lexer.TokenOperators["+"], Text: "+"},
				{Type: lexer.TokenText["identifier"], Text: "b"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "result"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenText["identifier"], Text: "add"},
				{Type: lexer.TokenPunctuation["("], Text: "("},
				{Type: lexer.TokenLiterals["number"], Text: "5"},
				{Type: lexer.TokenPunctuation[","], Text: ","},
				{Type: lexer.TokenLiterals["number"], Text: "10"},
				{Type: lexer.TokenPunctuation[")"], Text: ")"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
			},
		},
		// Test for loop and conditionals
		{
			Input: `
			for (var i = 0; i < 10; i++) {
				if (i % 2 == 0) {
					continue;
				} else {
					break;
				}
			}
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["for"], Text: "for"},
				{Type: lexer.TokenPunctuation["("], Text: "("},
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "i"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenLiterals["number"], Text: "0"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenText["identifier"], Text: "i"},
				{Type: lexer.TokenOperators["<"], Text: "<"},
				{Type: lexer.TokenLiterals["number"], Text: "10"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenText["identifier"], Text: "i"},
				{Type: lexer.TokenOperators["++"], Text: "++"},
				{Type: lexer.TokenPunctuation[")"], Text: ")"},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenKeywords["if"], Text: "if"},
				{Type: lexer.TokenPunctuation["("], Text: "("},
				{Type: lexer.TokenText["identifier"], Text: "i"},
				{Type: lexer.TokenOperators["%"], Text: "%"},
				{Type: lexer.TokenLiterals["number"], Text: "2"},
				{Type: lexer.TokenOperators["=="], Text: "=="},
				{Type: lexer.TokenLiterals["number"], Text: "0"},
				{Type: lexer.TokenPunctuation[")"], Text: ")"},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenKeywords["continue"], Text: "continue"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
				{Type: lexer.TokenKeywords["else"], Text: "else"},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenKeywords["break"], Text: "break"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
			},
		},
		// Test string literals and escape sequences
		{
			Input: `
			var str = "Hello, world!\n";
			var escapedStr = "This is a \"quoted\" string.";
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "str"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenLiterals["string"], Text: "Hello, world!\n"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "escapedStr"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenLiterals["string"], Text: "This is a \"quoted\" string."},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
			},
		},
		// Test boolean literals and logical operators
		{
			Input: `
			var a = true && false || !true;
			var b = !!a;
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "a"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenKeywords["true"], Text: "true"},
				{Type: lexer.TokenOperators["&&"], Text: "&&"},
				{Type: lexer.TokenKeywords["false"], Text: "false"},
				{Type: lexer.TokenOperators["||"], Text: "||"},
				{Type: lexer.TokenOperators["!"], Text: "!"},
				{Type: lexer.TokenKeywords["true"], Text: "true"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "b"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenOperators["!"], Text: "!"},
				{Type: lexer.TokenOperators["!"], Text: "!"},
				{Type: lexer.TokenText["identifier"], Text: "a"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
			},
		},
		// Test object literals
		{
			Input: `
			var obj = {
				key1: "value1",
				key2: 42,
				key3: true
			};
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "obj"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenText["identifier"], Text: "key1"},
				{Type: lexer.TokenPunctuation[":"], Text: ":"},
				{Type: lexer.TokenLiterals["string"], Text: "value1"},
				{Type: lexer.TokenPunctuation[","], Text: ","},
				{Type: lexer.TokenText["identifier"], Text: "key2"},
				{Type: lexer.TokenPunctuation[":"], Text: ":"},
				{Type: lexer.TokenLiterals["number"], Text: "42"},
				{Type: lexer.TokenPunctuation[","], Text: ","},
				{Type: lexer.TokenText["identifier"], Text: "key3"},
				{Type: lexer.TokenPunctuation[":"], Text: ":"},
				{Type: lexer.TokenKeywords["true"], Text: "true"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
			},
		},
		// Test array literals and access
		{
			Input: `
			var arr = [1, 2, 3];
			var first = arr[0];
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "arr"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenPunctuation["["], Text: "["},
				{Type: lexer.TokenLiterals["number"], Text: "1"},
				{Type: lexer.TokenPunctuation[","], Text: ","},
				{Type: lexer.TokenLiterals["number"], Text: "2"},
				{Type: lexer.TokenPunctuation[","], Text: ","},
				{Type: lexer.TokenLiterals["number"], Text: "3"},
				{Type: lexer.TokenPunctuation["]"], Text: "]"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "first"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenText["identifier"], Text: "arr"},
				{Type: lexer.TokenPunctuation["["], Text: "["},
				{Type: lexer.TokenLiterals["number"], Text: "0"},
				{Type: lexer.TokenPunctuation["]"], Text: "]"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
			},
		},
		// Test class declaration
		{
			Input: `
			class Person {
				constructor(name) {
					this.name = name;
				}
				getName() {
					return this.name;
				}
			}
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["class"], Text: "class"},
				{Type: lexer.TokenText["identifier"], Text: "Person"},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenText["identifier"], Text: "constructor"},
				{Type: lexer.TokenPunctuation["("], Text: "("},
				{Type: lexer.TokenText["identifier"], Text: "name"},
				{Type: lexer.TokenPunctuation[")"], Text: ")"},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenText["identifier"], Text: "this"},
				{Type: lexer.TokenPunctuation["."], Text: "."},
				{Type: lexer.TokenText["identifier"], Text: "name"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenText["identifier"], Text: "name"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
				{Type: lexer.TokenText["identifier"], Text: "getName"},
				{Type: lexer.TokenPunctuation["("], Text: "("},
				{Type: lexer.TokenPunctuation[")"], Text: ")"},
				{Type: lexer.TokenPunctuation["{"], Text: "{"},
				{Type: lexer.TokenKeywords["return"], Text: "return"},
				{Type: lexer.TokenText["identifier"], Text: "this"},
				{Type: lexer.TokenPunctuation["."], Text: "."},
				{Type: lexer.TokenText["identifier"], Text: "name"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
				{Type: lexer.TokenPunctuation["}"], Text: "}"},
			},
		},
		// Test regular expressions
		{
			Input: `
			var regex = /ab+c/;
			`,
			Expected: []lexer.GojoToken{
				{Type: lexer.TokenKeywords["var"], Text: "var"},
				{Type: lexer.TokenText["identifier"], Text: "regex"},
				{Type: lexer.TokenOperators["="], Text: "="},
				{Type: lexer.TokenLiterals["regexp"], Text: "/ab+c/"},
				{Type: lexer.TokenPunctuation[";"], Text: ";"},
			},
		},
	}

	lastTestFailed := false
	for i, test := range tests {
		l := lexer.New(test.Input)
		for _, expectedToken := range test.Expected {
			tok := l.NextToken()
			if env.Verbose {
				fmt.Printf("Expected: %v, Got: %v\n", expectedToken, tok)
			}
			if tok.Type != expectedToken.Type || tok.Text != expectedToken.Text {
				fmt.Printf("❌ Test %d failed. Expected: %v, Got: %v\n", i+1, expectedToken, tok)
				lastTestFailed = true
				break
			}
		}
		if lastTestFailed {
			lastTestFailed = false
		} else {
			fmt.Printf("✅ Test %d passed.\n", i+1)
		}
	}
}
