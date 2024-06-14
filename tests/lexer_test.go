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
		var token GojoToken = lexer.NextToken()
		if token.Type != expectedToken.Type || token.Text != expectedToken.Text {
			t.Errorf("Token %2d: \nExpected: %v\nReceived: %v\n", i, expectedToken, token)
		} else {
			t.Logf("Token %2d: %v", i, token)
		}
	}
}

func NewToken(tokenTypeStr string, text ...string) GojoToken {
	tokenMaps := []map[string]*GojoTokenType{
		TokenKeywords,
		TokenPunctuation,
		TokenOperators,
		TokenText,
		TokenLiterals,
	}
	var tokenType *GojoTokenType
	var tokenExists bool = false

	for _, tokenMap := range tokenMaps {
		tokenType, tokenExists = tokenMap[tokenTypeStr]
		if tokenExists {
			tokenText := tokenTypeStr
			if len(text) > 0 {
				tokenText = text[0]
			}
			return GojoToken{
				Text: tokenText,
				Type: tokenType,
			}
		}
	}
	panic("Token not found")
}

func NewNumber(number string) GojoToken {
	return NewToken("number", number)
}

func NewID(id string) GojoToken {
	return NewToken("identifier", id)
}

func NewString(str string) GojoToken {
	return NewToken("string", str)
}

var lexerTestCases = []LexerTestCase{
	// Test variable declarations and basic arithmetic operations
	{
		Name: "VariablesAndArithmetic",
		Expected: []GojoToken{
			NewToken("var"), NewID("x"), NewToken("="), NewNumber("5"), NewToken(";"),
			NewToken("var"), NewID("y"), NewToken("="), NewNumber("10"), NewToken(";"),
			NewToken("var"), NewID("z"), NewToken("="), NewID("x"), NewToken("+"), NewID("y"), NewToken(";"),
		},
	},
	{
		Name: "Functions",
		Expected: []GojoToken{
			NewToken("function"), NewID("add"), NewToken("("), NewID("a"), NewToken(","), NewID("b"), NewToken(")"), NewToken("{"),
			NewToken("return"), NewID("a"), NewToken("+"), NewID("b"), NewToken(";"),
			NewToken("}"),
			NewToken("var"), NewID("result"), NewToken("="), NewID("add"), NewToken("("), NewNumber("5"), NewToken(","), NewNumber("10"), NewToken(")"), NewToken(";"),
		},
	},
	{
		Name: "LoopsAndConditionals",
		Expected: []GojoToken{
			NewToken("for"), NewToken("("), NewToken("var"), NewID("i"), NewToken("="), NewNumber("0"), NewToken(";"),
			NewID("i"), NewToken("<"), NewNumber("10"), NewToken(";"), NewID("i"), NewToken("++"), NewToken(")"), NewToken("{"),
			NewToken("if"), NewToken("("), NewID("i"), NewToken("%"), NewNumber("2"), NewToken("=="), NewNumber("0"), NewToken(")"), NewToken("{"),
			NewToken("continue"), NewToken(";"),
			NewToken("}"), NewToken("else"), NewToken("{"),
			NewToken("break"), NewToken(";"),
			NewToken("}"),
			NewToken("}"),
		},
	},
	{
		Name: "Strings",
		Expected: []GojoToken{
			NewToken("var"), NewID("str"), NewToken("="), NewString("Hello, world!\n"), NewToken(";"),
			NewToken("var"), NewID("escapedStr"), NewToken("="), NewString("This is a \"quoted\" string."), NewToken(";"),
		},
	},
	{
		Name: "Booleans",
		Expected: []GojoToken{
			NewToken("var"), NewID("a"), NewToken("="), NewToken("true"), NewToken("&&"), NewToken("false"), NewToken("||"), NewToken("!"), NewToken("true"), NewToken(";"),
			NewToken("var"), NewID("b"), NewToken("="), NewToken("!"), NewToken("!"), NewID("a"), NewToken(";"),
		},
	},
	{
		Name: "Objects",
		Expected: []GojoToken{
			NewToken("var"), NewID("obj"), NewToken("="), NewToken("{"),
			NewID("key1"), NewToken(":"), NewString("value1"), NewToken(","),
			NewID("key2"), NewToken(":"), NewNumber("42"), NewToken(","),
			NewID("key3"), NewToken(":"), NewToken("true"),
			NewToken("}"), NewToken(";"),
		},
	},
	{
		Name: "Arrays",
		Expected: []GojoToken{
			NewToken("var"), NewID("arr"), NewToken("="), NewToken("["), NewNumber("1"), NewToken(","), NewNumber("2"), NewToken(","), NewNumber("3"), NewToken("]"), NewToken(";"),
			NewToken("var"), NewID("first"), NewToken("="), NewID("arr"), NewToken("["), NewNumber("0"), NewToken("]"), NewToken(";"),
		},
	},
	{
		Name: "ClassDeclaration",
		Expected: []GojoToken{
			NewToken("class"), NewID("Person"), NewToken("{"),
			NewID("constructor"), NewToken("("), NewID("name"), NewToken(")"), NewToken("{"),
			NewToken("this"), NewToken("."), NewID("name"), NewToken("="), NewID("name"), NewToken(";"),
			NewToken("}"),
			NewID("getName"), NewToken("("), NewToken(")"), NewToken("{"),
			NewToken("return"), NewToken("this"), NewToken("."), NewID("name"), NewToken(";"),
			NewToken("}"),
			NewToken("}"),
		},
	},
	{
		Name: "Regexp",
		Expected: []GojoToken{
			NewToken("var"), NewID("regex"), NewToken("="), NewToken("regexp", "/ab+c/"), NewToken(";"),
		},
	},
	{
		Name: "MultiCharacterOperators",
		Expected: []GojoToken{
			NewToken("var"), NewID("a"), NewToken("="), NewNumber("2"), NewToken("++"), NewToken(";"),
			NewToken("var"), NewID("b"), NewToken("="), NewToken("++"), NewNumber("16"), NewToken(";"),
			NewID("a"), NewToken("+="), NewNumber("5"), NewToken(";"),
			NewToken("var"), NewID("c"), NewToken("="), NewID("a"), NewToken("=="), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("d"), NewToken("="), NewID("a"), NewToken("!="), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("e"), NewToken("="), NewID("a"), NewToken("!=="), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("f"), NewToken("="), NewID("a"), NewToken("==="), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("g"), NewToken("="), NewID("a"), NewToken("<="), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("h"), NewToken("="), NewID("a"), NewToken(">="), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("i"), NewToken("="), NewID("a"), NewToken("&&"), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("j"), NewToken("="), NewID("a"), NewToken("||"), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("k"), NewToken("="), NewID("a"), NewToken("<<"), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("l"), NewToken("="), NewID("a"), NewToken(">>"), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("m"), NewToken("="), NewID("a"), NewToken(">>>"), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("n"), NewToken("="), NewID("a"), NewToken("**"), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("o"), NewToken("="), NewID("a"), NewToken("??"), NewID("b"), NewToken(";"),
			NewToken("var"), NewID("p"), NewToken("="), NewID("a"), NewToken("?."), NewID("b"), NewToken(";"),
		},
	},
}
