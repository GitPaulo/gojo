package parser

import (
	"fmt"
	"gojo/config"
	"gojo/lexer"
	"strconv"
)

const (
	LOWEST  = iota
	SUM     // +
	PRODUCT // *, /
)

type Parser struct {
	l             *lexer.Lexer
	errors        []string
	curToken      lexer.GojoToken
	curTokenStart int
	curTokenEnd   int
	peekToken     lexer.GojoToken
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	var token = p.peekToken
	p.curToken = token
	p.curTokenStart = p.l.Start
	p.curTokenEnd = p.l.End
	p.peekToken = p.l.NextToken()

	if config.LoadConfig().Verbose {
		fmt.Println("---")
		fmt.Println("Current Token:", p.curToken)
		fmt.Println("Start Position:", p.curTokenStart)
		fmt.Println("End Position:", p.curTokenEnd)
	}
}

/**
 * Parsing functions
 */

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type.Label != "eof" {
		stmt := p.parseStatement()
		if stmt != nil {
			if config.LoadConfig().Verbose {
				fmt.Println("Generated Statement:", stmt)
			}
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type.Label {
	case "var", "let", "const":
		return p.parseVariableDeclarationStatement()
	default:
		return nil
	}
}

func (p *Parser) parseVariableDeclarationStatement() *VariableDeclaration {
	stmt := &VariableDeclaration{Token: p.curToken}

	if !p.expectPeek("identifier") {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Text}

	if !p.expectPeek("=") {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(";") {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	leftExp := p.parseAtomicExpression()

	for !p.peekTokenIs(";") && precedence < p.peekPrecedence() {
		p.nextToken()
		leftExp = p.parseInfixExpression(leftExp)
	}

	return leftExp
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	exp := &BinaryExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Text,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(precedence)
	return exp
}

func (p *Parser) curPrecedence() int {
	switch p.curToken.Type.Label {
	case "+", "-":
		return SUM
	case "*", "/":
		return PRODUCT
	default:
		return LOWEST
	}
}

func (p *Parser) peekPrecedence() int {
	switch p.peekToken.Type.Label {
	case "+", "-":
		return SUM
	case "*", "/":
		return PRODUCT
	default:
		return LOWEST
	}
}

func (p *Parser) parseAtomicExpression() Expression {
	switch p.curToken.Type.Label {
	case "number":
		return p.parseIntegerLiteral()
	case "identifier":
		return p.parseIdentifier()
	case "(":
		return p.parseGroupedExpression()
	default:
		return nil
	}
}

func (p *Parser) parseIntegerLiteral() *IntegerLiteral {
	lit := &IntegerLiteral{Token: p.curToken}
	lit.Value, _ = strconv.ParseInt(p.curToken.Text, 0, 64)
	return lit
}

func (p *Parser) parseIdentifier() *Identifier {
	return &Identifier{Token: p.curToken, Value: p.curToken.Text}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken() // Consume "("
	exp := p.parseExpression(LOWEST)
	p.nextToken() // Consume ")"
	return exp
}

/**
 * Helper functions
 */

func (p *Parser) expectPeek(t string) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) peekTokenIs(t string) bool {
	return p.peekToken.Type.Label == t
}

/**
 * Error handling
 */

func (p *Parser) Errors() []string {
	return p.errors
}
