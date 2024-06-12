package lexer

import (
	"fmt"
	"os"
	"unicode"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	nextPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	// Exported
	Start int // start position of the current token
	End   int // end position of the current token
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func newToken(tokenType GojoTokenType, text string) GojoToken {
	return GojoToken{
		Text: text,
		Type: tokenType,
	}
}

func (l *Lexer) NextToken() GojoToken {
	var tok GojoToken

	l.skipWhitespace()

	switch l.ch {
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
			tok = newToken(TokenOperators["/"], string(l.ch))
		}
	case '=':
		tok = l.readTwoCharOperator('=', TokenOperators["="], TokenOperators["=="], TokenOperators["==="])
	case '+':
		tok = l.readTwoCharOperator('+', TokenOperators["+"], TokenOperators["+="], TokenOperators["++"])
	case '-':
		tok = l.readTwoCharOperator('-', TokenOperators["-"], TokenOperators["-="], TokenOperators["--"])
	case '*':
		tok = l.readTwoCharOperator('*', TokenOperators["*"], TokenOperators["*="], TokenOperators["**"])
	case '!':
		tok = l.readTwoCharOperator('=', TokenOperators["!"], TokenOperators["!="], TokenOperators["!=="])
	case '<':
		tok = l.readTwoCharOperator('<', TokenOperators["<"], TokenOperators["<="], TokenOperators["<<"])
	case '>':
		tok = l.readTwoCharOperator('>', TokenOperators[">"], TokenOperators[">="], TokenOperators[">>"])
	case '&':
		tok = l.readTwoCharOperator('&', TokenOperators["&"], TokenOperators["&="], TokenOperators["&&"])
	case '|':
		tok = l.readTwoCharOperator('|', TokenOperators["|"], TokenOperators["|="], TokenOperators["||"])
	case '^':
		tok = l.readTwoCharOperator('^', TokenOperators["^"], TokenOperators["^="], TokenOperators["^"])
	case '%':
		tok = l.readTwoCharOperator('%', TokenOperators["%"], TokenOperators["%="], TokenOperators["%"])
	case '.':
		tok = newToken(TokenPunctuation["."], string(l.ch))
		if l.peekChar() == '.' && l.peekCharTwo() == '.' {
			l.readChar()
			l.readChar()
			tok = newToken(TokenPunctuation["..."], "...")
		}
	case ',':
		tok = newToken(TokenPunctuation[","], string(l.ch))
	case ';':
		tok = newToken(TokenPunctuation[";"], string(l.ch))
	case ':':
		tok = newToken(TokenPunctuation[":"], string(l.ch))
	case '(':
		tok = newToken(TokenPunctuation["("], string(l.ch))
	case ')':
		tok = newToken(TokenPunctuation[")"], string(l.ch))
	case '{':
		tok = newToken(TokenPunctuation["{"], string(l.ch))
	case '}':
		tok = newToken(TokenPunctuation["}"], string(l.ch))
	case '[':
		tok = newToken(TokenPunctuation["["], string(l.ch))
	case ']':
		tok = newToken(TokenPunctuation["]"], string(l.ch))
	case '?':
		tok = newToken(TokenPunctuation["?"], string(l.ch))
	case '"', '\'', '`': // Handle strings with all three quote types
		tok = l.readString(l.ch)
	case 0:
		tok = GojoToken{Type: TokenText["eof"], Text: ""}
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			tokenType, ok := TokenKeywords[identifier]
			if !ok {
				tokenType = TokenText["identifier"]
			}
			return GojoToken{Text: identifier, Type: tokenType}
		} else if isDigit(l.ch) {
			number := l.readNumber()
			return GojoToken{Text: number, Type: TokenLiterals["number"]}
		} else {
			fmt.Println("Unknown token: ", string(l.ch))
			os.Exit(1)
		}
	}

	l.readChar()
	return tok
}

/**
 * Read methods
 */

func (l *Lexer) readTwoCharOperator(expected byte, single, double, triple GojoTokenType) GojoToken {
	ch := string(l.ch)
	if l.peekChar() == expected {
		l.readChar()
		ch += string(l.ch)
		if l.peekChar() == expected && (single.Label == "==" || single.Label == "!=" || single.Label == "<" || single.Label == ">" || single.Label == "+" || single.Label == "-" || single.Label == "*" || single.Label == "/" || single.Label == "%" || single.Label == "&" || single.Label == "|" || single.Label == "^") {
			l.readChar()
			ch += string(l.ch)
			return newToken(triple, ch)
		}
		return newToken(double, ch)
	}
	return newToken(single, ch)
}

func (l *Lexer) readString(quoteType byte) GojoToken {
	l.readChar() // Consume the opening quote

	var text string
	for {
		ch := l.ch
		if ch == 0 {
			fmt.Println("Unterminated string literal")
			os.Exit(1)
		} else if ch == quoteType {
			l.readChar() // Consume closing quote
			break
		} else if ch == '\\' {
			l.readChar() // Consume escape character
			switch l.ch {
			case 'n':
				text += "\n"
			case 't':
				text += "\t"
			case 'r':
				text += "\r"
			case '\\':
				text += "\\"
			case quoteType: // Handle escaped quote of the same type
				text += string(quoteType)
			case '\'':
				text += "'"
			case '"':
				text += "\""
			case '`':
				text += "`"
			case 'x':
				hex := l.readHex(2)
				text += fmt.Sprintf("\\x%s", hex)
			case 'u':
				hex := l.readHex(4)
				text += fmt.Sprintf("\\u%s", hex)
			default:
				fmt.Println("Invalid escape sequence in string: \\", string(l.ch))
				os.Exit(1)
			}
			l.readChar() // Consume escaped character
		} else {
			text += string(ch)
			l.readChar()
		}
	}

	return GojoToken{Type: TokenLiterals["string"], Text: text}
}

func (l *Lexer) readRegex() GojoToken {
	startPos := l.position
	for {
		l.readChar()
		if l.ch == '/' && l.input[l.position-1] != '\\' {
			break
		}
		if l.ch == 0 {
			fmt.Println("Unterminated regex literal")
			os.Exit(1)
		}
	}
	l.readChar() // Move past the closing '/'
	return GojoToken{Type: TokenLiterals["regexp"], Text: l.input[startPos:l.position]}
}

func (l *Lexer) readHex(length int) string {
	var hex string
	for i := 0; i < length; i++ {
		l.readChar()
		hex += string(l.ch)
	}
	return hex
}

func (l *Lexer) readChar() {
	l.ch = l.peekChar()
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
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

/**
 * Character checks
 */

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

/**
 * Skips: Whitespace and comments
 */

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		l.readChar()
	}
}

func (l *Lexer) skipInlineComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	l.readChar() // consume the '*'
	for {
		l.readChar()
		if l.ch == '*' && l.peekChar() == '/' {
			l.readChar() // consume the '*'
			l.readChar() // consume the '/'
			break
		}
		if l.ch == 0 {
			break // End of input reached before the end of the block comment
		}
	}
}
