package lexer

import (
	"fmt"
	"gojo/config"
	"os"
	"unicode"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	nextPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
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
	var tok GojoToken

	l.skipWhitespace()

	if config.LoadConfig().MegaVerbose {
		fmt.Printf("Current character: %c\n", l.ch)
	}

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
			tok = l.NewToken(TokenOperators["/"], string(l.ch))
		}
	case '=':
		tok = l.readMultiCharOperator(TokenOperators["="], TokenOperators["=="], TokenOperators["==="])
	case '+':
		tok = l.readMultiCharOperator(TokenOperators["+"], TokenOperators["+="], TokenOperators["++"])
	case '-':
		tok = l.readMultiCharOperator(TokenOperators["-"], TokenOperators["-="], TokenOperators["--"])
	case '*':
		tok = l.readMultiCharOperator(TokenOperators["*"], TokenOperators["*="], TokenOperators["**"])
	case '!':
		tok = l.readMultiCharOperator(TokenOperators["!"], TokenOperators["!="], TokenOperators["!=="])
	case '<':
		tok = l.readMultiCharOperator(TokenOperators["<"], TokenOperators["<="], TokenOperators["<<"])
	case '>':
		tok = l.readMultiCharOperator(TokenOperators[">"], TokenOperators[">="], TokenOperators[">>"])
	case '&':
		tok = l.readMultiCharOperator(TokenOperators["&"], TokenOperators["&&"])
	case '|':
		tok = l.readMultiCharOperator(TokenOperators["|"], TokenOperators["|="], TokenOperators["||"])
	case '^':
		tok = l.readMultiCharOperator(TokenOperators["^"], TokenOperators["^="], TokenOperators["^^"])
	case '%':
		tok = l.readMultiCharOperator(TokenOperators["%"], TokenOperators["%="])
	case '.':
		tok = l.NewToken(TokenPunctuation["."], string(l.ch))
		if l.peekChar() == '.' && l.peekCharTwo() == '.' {
			l.readChar()
			l.readChar()
			tok = l.NewToken(TokenPunctuation["..."], "...")
		}
	case ',':
		tok = l.NewToken(TokenPunctuation[","], string(l.ch))
	case ';':
		tok = l.NewToken(TokenPunctuation[";"], string(l.ch))
	case ':':
		tok = l.NewToken(TokenPunctuation[":"], string(l.ch))
	case '(':
		tok = l.NewToken(TokenPunctuation["("], string(l.ch))
	case ')':
		tok = l.NewToken(TokenPunctuation[")"], string(l.ch))
	case '{':
		tok = l.NewToken(TokenPunctuation["{"], string(l.ch))
	case '}':
		tok = l.NewToken(TokenPunctuation["}"], string(l.ch))
	case '[':
		tok = l.NewToken(TokenPunctuation["["], string(l.ch))
	case ']':
		tok = l.NewToken(TokenPunctuation["]"], string(l.ch))
	case '?':
		tok = l.NewToken(TokenPunctuation["?"], string(l.ch))
	case '"', '\'', '`': // Handle strings with all three quote types
		return l.readString(l.ch)
	case 0:
		tok = GojoToken{Type: TokenText["eof"], Text: ""}
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			tokenType, ok := TokenKeywords[identifier]
			if !ok {
				tokenType = TokenText["identifier"]
			}
			return l.NewToken(tokenType, identifier)
		} else if isDigit(l.ch) {
			number := l.readNumber()
			return l.NewToken(TokenLiterals["number"], number)
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

func (l *Lexer) readMultiCharOperator(options ...*GojoTokenType) GojoToken {
	ch := string(l.ch)
	for _, option := range options {
		if option == nil {
			continue
		} else if len(option.Label) < 2 {
			continue
		} else if l.peekChar() == option.Label[1] {
			l.readChar()
			ch += string(l.ch)
			return l.NewToken(option, ch)
		}
	}
	return l.NewToken(options[0], ch)
}

func (l *Lexer) readString(quoteType byte) GojoToken {
	l.readChar() // Consume the opening quote

	var text string
	for {
		ch := l.ch
		if ch == 0 {
			panic("Unterminated string literal")
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

	return l.NewToken(TokenLiterals["string"], text)
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
	return l.NewToken(TokenLiterals["regexp"], l.input[startPos:l.position])
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
	if config.LoadConfig().MegaVerbose {
		fmt.Println("Reading character: ", string(l.ch))
	}
	if l.ch == '\n' {
		if config.LoadConfig().Verbose {
			fmt.Println("░ Newline detected")
		}
		l.Line++
	}
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
	startPos := l.position
	hasDecimalPoint := false
	isScientificNotation := false

	for {
		if isDigit(l.ch) {
			// continue reading digits
		} else if l.ch == '.' && !hasDecimalPoint {
			// first decimal point is valid
			hasDecimalPoint = true
		} else if (l.ch == 'e' || l.ch == 'E') && !isScientificNotation {
			// scientific notation (e.g., 1e10)
			isScientificNotation = true
			l.readChar()
			if l.ch == '+' || l.ch == '-' {
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
