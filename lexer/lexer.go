package lexer

import (
	"fmt"
	"gojo/config"
	"os"
	"strings"
	"unicode"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	nextPosition int  // current reading position in input (after current char)
	curChar      byte // current char under examination
	// Exported
	Line  int // current line number
	Start int // start position of the current token
	End   int // end position of the current token
}

func New(input string) *Lexer {
	l := &Lexer{input: input, Line: 1}
	l.readChar()
	return l
}

func (l *Lexer) NewToken(tokenType *GojoTokenType, text string) GojoToken {
	return GojoToken{
		Text: text,
		Type: tokenType,
		Line: l.Line,
	}
}

func (l *Lexer) NextToken() GojoToken {
	var token GojoToken

	l.skipWhitespace()

	if config.LoadConfig().MegaVerbose {
		fmt.Printf("Current character: %c\n", l.curChar)
	}

	switch l.curChar {
	case '/':
		if l.peekChar() == '/' {
			l.skipInlineComment()
			return l.NextToken()
		} else if l.peekChar() == '*' {
			l.skipBlockComment()
			return l.NextToken()
		} else if l.peekChar() == '/' || l.peekChar() == '*' || l.peekChar() == '+' || l.peekChar() == '-' {
			return l.readRegex()
		} else {
			token = l.NewToken(TokenOperators["/"], string(l.curChar))
		}
	case '=', '+', '-', '*', '!', '<', '>', '&', '|', '^', '%', '?':
		token = l.readOperator()
	case '.':
		token = l.NewToken(TokenPunctuation["."], string(l.curChar))
		if l.peekChar() == '.' && l.peekCharTwo() == '.' {
			l.readChar()
			l.readChar()
			token = l.NewToken(TokenPunctuation["..."], "...")
		}
	case ',':
		token = l.NewToken(TokenPunctuation[","], string(l.curChar))
	case ';':
		token = l.NewToken(TokenPunctuation[";"], string(l.curChar))
	case ':':
		token = l.NewToken(TokenPunctuation[":"], string(l.curChar))
	case '(':
		token = l.NewToken(TokenPunctuation["("], string(l.curChar))
	case ')':
		token = l.NewToken(TokenPunctuation[")"], string(l.curChar))
	case '{':
		token = l.NewToken(TokenPunctuation["{"], string(l.curChar))
	case '}':
		token = l.NewToken(TokenPunctuation["}"], string(l.curChar))
	case '[':
		token = l.NewToken(TokenPunctuation["["], string(l.curChar))
	case ']':
		token = l.NewToken(TokenPunctuation["]"], string(l.curChar))
	case '"', '\'', '`': // Handle strings with all three quote types
		return l.readString(l.curChar)
	case 0:
		token = GojoToken{Type: TokenText["eof"], Text: ""}
	default:
		// Note: letters can be a lot! (e.g., keywords, literals and identifiers)
		if isLetter(l.curChar) {
			word := l.readWord()
			tokenType, ok := TokenKeywords[word]
			if !ok {
				tokenType, ok = TokenLiterals[word]
				if !ok {
					tokenType = TokenText["identifier"]
				}
			}
			return l.NewToken(tokenType, word)
		} else if isDigit(l.curChar) {
			number := l.readNumber()
			return l.NewToken(TokenLiterals["number"], number)
		} else {
			fmt.Println("Unknown token: ", string(l.curChar))
			os.Exit(1)
		}
	}

	l.readChar()

	return token
}

/**
 * Read methods
 */

func (l *Lexer) readOperator() GojoToken {
	operatorStr := string(l.curChar)
	tokenType, validToken := TokenOperators[operatorStr]

	if validToken {
		var testOperator strings.Builder
		testOperator.WriteString(operatorStr)

		for {
			// Check if the next character is a valid operator
			nextChar := l.peekChar()
			testOperator.WriteByte(nextChar)
			testString := testOperator.String()
			testType, testValid := TokenOperators[testString]

			if !testValid {
				break // Not part of a valid operator, break early
			}

			tokenType = testType
			operatorStr = testString
			l.readChar()
		}
	}

	if !validToken {
		panic("Invalid operator read: " + operatorStr)
	}

	return l.NewToken(tokenType, operatorStr)
}

func (l *Lexer) readString(quoteType byte) GojoToken {
	l.readChar() // Consume the opening quote
	var text string

	for {
		readChar := l.curChar
		if readChar == 0 {
			panic("Unterminated string literal: " + text)
		} else if readChar == quoteType {
			l.readChar() // Consume closing quote
			break
		} else if readChar == '\\' {
			text += l.readEscapeSequence(quoteType)
		} else {
			text += string(readChar)
			l.readChar()
		}
	}

	return l.NewToken(TokenLiterals["string"], text)
}

func (l *Lexer) readEscapeSequence(quoteType byte) string {
	l.readChar() // Consume escape character
	switch l.curChar {
	case 'n':
		return "\n"
	case 't':
		return "\t"
	case 'r':
		return "\r"
	case '\\':
		return "\\"
	case '\'':
		return "'"
	case '"':
		return "\""
	case '`':
		return "`"
	case 'x':
		hex := l.readHex(2)
		return fmt.Sprintf("\\x%s", hex)
	case 'u':
		hex := l.readHex(4)
		return fmt.Sprintf("\\u%s", hex)
	case quoteType:
		return string(quoteType)
	default:
		fmt.Println("Invalid escape sequence in string: \\", string(l.curChar))
		os.Exit(1)
		return ""
	}
}

func (l *Lexer) readRegex() GojoToken {
	startPos := l.position
	for {
		l.readChar()
		if l.curChar == '/' && l.input[l.position-1] != '\\' {
			break
		}
		if l.curChar == 0 {
			fmt.Println("Unterminated regex literal")
			os.Exit(1)
		}
	}
	l.readChar() // Move past the closing '/'
	return l.NewToken(TokenLiterals["regexp"], l.input[startPos:l.position])
}

func (l *Lexer) readHex(length int) string {
	var hex string
	for i := 0; i < length; i++ {
		l.readChar()
		hex += string(l.curChar)
	}
	return hex
}

func (l *Lexer) readChar() {
	if config.LoadConfig().MegaVerbose {
		fmt.Println("Reading character: ", string(l.curChar))
	}
	if l.curChar == '\n' {
		if config.LoadConfig().Verbose {
			fmt.Println("░ Newline detected")
		}
		l.Line++
	}
	l.curChar = l.peekChar()
	l.Start = l.position
	if l.position < len(l.input) {
		l.position = l.nextPosition
		l.End = l.position
		l.nextPosition++
	} else {
		l.End = l.position
	}
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition]
}

func (l *Lexer) peekCharTwo() byte {
	if l.nextPosition+1 >= len(l.input) {
		return 0
	}
	return l.input[l.nextPosition+1]
}

func (l *Lexer) readNumber() string {
	startPos := l.position
	hasDecimalPoint := false
	isScientificNotation := false

	for {
		if isDigit(l.curChar) {
			// continue reading digits
		} else if l.curChar == '.' && !hasDecimalPoint {
			// first decimal point is valid
			hasDecimalPoint = true
		} else if (l.curChar == 'e' || l.curChar == 'E') && !isScientificNotation {
			// scientific notation (e.g., 1e10)
			isScientificNotation = true
			l.readChar()
			if l.curChar == '+' || l.curChar == '-' {
				l.readChar()
			}
			continue
		} else {
			break // break on non-digit character
		}
		l.readChar()
	}

	return l.input[startPos:l.position]
}

func (l *Lexer) readWord() string {
	pos := l.position
	for isLetter(l.curChar) || isDigit(l.curChar) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

/**
 * Character checks
 */

func isLetter(char byte) bool {
	return unicode.IsLetter(rune(char)) || char == '_'
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

/**
 * Skips: Whitespace and comments
 */

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.curChar)) {
		l.readChar()
	}
}

func (l *Lexer) skipInlineComment() {
	for l.curChar != '\n' && l.curChar != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	l.readChar() // consume the '*'
	for {
		l.readChar()
		if l.curChar == '*' && l.peekChar() == '/' {
			l.readChar() // consume the '*'
			l.readChar() // consume the '/'
			break
		}
		if l.curChar == 0 {
			break // End of input reached before the end of the block comment
		}
	}
}
