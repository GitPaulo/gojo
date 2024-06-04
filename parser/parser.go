package parser

import (
	"gojo/lexer"
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  lexer.GojoToken
	peekToken lexer.GojoToken
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type.Label != "eof" {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type.Label {
	case "var":
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek("identifier") {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Text}

	if !p.expectPeek("=") {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression()

	if p.peekTokenIs(";") {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression() Expression {
	switch p.curToken.Type.Label {
	case "number":
		return p.parseIntegerLiteral()
	default:
		return nil
	}
}

func (p *Parser) parseIntegerLiteral() *IntegerLiteral {
	lit := &IntegerLiteral{Token: p.curToken}

	return lit
}

func (p *Parser) expectPeek(t string) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) peekTokenIs(t string) bool {
	return p.peekToken.Type.Label == t
}

func (p *Parser) curTokenIs(t string) bool {
	return p.curToken.Type.Label == t
}

func (p *Parser) Errors() []string {
	return p.errors
}
