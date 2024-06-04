package lexer

import (
	"fmt"
	"os"
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

func newToken(tokenType GojoTokenType, ch byte) GojoToken {
	return GojoToken{
		Text: string(ch),
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
		} else {
			tok = newToken(tokenOperators["/"], l.ch)
		}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = newToken(tokenOperators["==="], l.ch)
			} else {
				tok = newToken(tokenOperators["=="], l.ch)
			}
		} else {
			tok = newToken(tokenOperators["="], l.ch)
		}
	case '+':
		tok = newToken(tokenOperators["+"], l.ch)
	case '-':
		tok = newToken(tokenOperators["-"], l.ch)
	case '*':
		tok = newToken(tokenOperators["*"], l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				tok = newToken(tokenOperators["!=="], l.ch)
			} else {
				tok = newToken(tokenOperators["!="], l.ch)
			}
		} else {
			tok = newToken(tokenOperators["!"], l.ch)
		}
	case '(':
		tok = newToken(tokenPunctuation["("], l.ch)
	case ')':
		tok = newToken(tokenPunctuation[")"], l.ch)
	case '{':
		tok = newToken(tokenPunctuation["{"], l.ch)
	case '}':
		tok = newToken(tokenPunctuation["}"], l.ch)
	case '[':
		tok = newToken(tokenPunctuation["["], l.ch)
	case ']':
		tok = newToken(tokenPunctuation["]"], l.ch)
	case ',':
		tok = newToken(tokenPunctuation[","], l.ch)
	case ';':
		tok = newToken(tokenPunctuation[";"], l.ch)
	case ':':
		tok = newToken(tokenPunctuation[":"], l.ch)
	case '.':
		tok = newToken(tokenPunctuation["."], l.ch)
	case '?':
		tok = newToken(tokenPunctuation["?"], l.ch)
	case '"', '\'', '`': // Handle strings with all three quote types
		tok = l.readString(l.ch)
	case 0:
		tok = GojoToken{Type: tokenText["eof"], Text: ""}
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			tokenType, ok := tokenKeywords[identifier]
			if !ok {
				tokenType = tokenText["identifier"]
			}
			return GojoToken{Text: identifier, Type: tokenType}
		} else if isDigit(l.ch) {
			number := l.readNumber()
			return GojoToken{Text: number, Type: tokenLiterals["number"]}
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
			case '\'': // Handle escaped single quote within double-quoted or backtick strings
				if quoteType != '"' && quoteType != '`' {
					fmt.Println("Invalid escape sequence in string: \\'", string(l.ch))
					os.Exit(1)
				} else {
					text += "'"
				}
			case '"': // Handle escaped double quote within single-quoted or backtick strings
				if quoteType != '\'' && quoteType != '`' {
					fmt.Println("Invalid escape sequence in string: \\\"", string(l.ch))
					os.Exit(1)
				} else {
					text += "\""
				}
			case 'x': // Handle hex escapes TODO
				// Implement logic to read and convert hex digits
				// text += handleHexEscape()
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

	return GojoToken{Type: tokenText["string"], Text: text}
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

func (l *Lexer) readNumber() string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.position]
}

/**
 * Character checks
 */

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

/**
 * Skips: Whitespace and comments
 */

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
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
