package lexer

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	nextPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

/**
New
*/

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
	case 0:
		tok = GojoToken{Type: tokenText["eof"], Text: ""}
	default:
		if isLetter(l.ch) {
			identifier := l.readIdentifier()
			tokenType, ok := tokenKeywords[identifier]
			if !ok {
				tokenType = GojoTokenType{Label: "identifier"}
			}
			return GojoToken{Text: identifier, Type: tokenType}
		} else if isDigit(l.ch) {
			number := l.readNumber()
			return GojoToken{Text: number, Type: GojoTokenType{Label: "number"}}
		} else {
			tok = newToken(GojoTokenType{Label: "illegal"}, l.ch)
		}
	}

	l.readChar()
	return tok
}

/**
Read
- Use byte slices to read identifiers and numbers, converting to strings only when necessary.
*/

func (l *Lexer) readChar() {
	l.ch = l.peekChar()
	l.position = l.nextPosition
	l.nextPosition += 1
}

func (l *Lexer) readString() string {
	pos := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[pos:l.position]
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

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

/**
Is
*/

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

/**
Skips
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
