package parser

import (
	"fmt"
	"gojo/config"
	"gojo/lexer"
	"strconv"
)

const (
	LOWEST      = iota
	EQUALITY    // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type Parser struct {
	l             *lexer.Lexer
	errors        []string
	curLine       int
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
	p.curLine = p.l.Line

	if config.LoadConfig().Verbose {
		fmt.Println("╔═══ nextToken() ")
		fmt.Println("Current Token:", p.curToken)
		fmt.Println("Start Position:", p.curTokenStart)
		fmt.Println("End Position:", p.curTokenEnd)
		fmt.Println("Current Line:", p.curLine)
		fmt.Println("Peek Token:", p.peekToken)
	}

	p.peekToken = p.l.NextToken()
}

/**
 * Parsing functions
 */

func (p *Parser) ParseProgram() *Program {
	program := &Program{Start: 0}
	program.Statements = []Statement{}

	for p.curToken.Type.Label != "eof" {
		stmt := p.parseStatement()
		if stmt != nil {
			if config.LoadConfig().Verbose {
				fmt.Println("╚══ parseStatement():", stmt)
			}
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	program.End = p.curTokenEnd

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type.Label {
	case "var", "let", "const":
		return p.parseVariableDeclarationStatement()
	case "function":
		return p.parseFunctionDeclaration()
	case "if":
		return p.parseIfStatement()
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

func (p *Parser) parseFunctionDeclaration() *FunctionDeclaration {
	stmt := &FunctionDeclaration{Token: p.curToken}

	if !p.expectPeek("identifier") {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Text}

	if !p.expectPeek("(") {
		return nil
	}

	stmt.Parameters = p.parseFunctionParameters()

	if !p.expectPeek("{") {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseFunctionParameters() []*Identifier {
	var identifiers []*Identifier

	if p.peekTokenIs(")") {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &Identifier{Token: p.curToken, Value: p.curToken.Text}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(",") {
		p.nextToken()
		p.nextToken()
		ident := &Identifier{Token: p.curToken, Value: p.curToken.Text}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(")") {
		return nil
	}

	return identifiers
}

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{Token: p.curToken}

	if !p.expectPeek("(") {
		return nil
	}

	p.nextToken()
	stmt.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(")") {
		return nil
	}

	if !p.expectPeek("{") {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	if p.peekTokenIs("else") {
		p.nextToken()

		if !p.expectPeek("{") {
			return nil
		}

		stmt.Alternative = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Token: p.curToken}
	block.Statements = []Statement{}

	p.nextToken()

	for !p.curTokenIs("}") && p.curToken.Type.Label != "eof" {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
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
	case "==", "!=", "<", ">", "<=", ">=":
		return EQUALITY
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
	case "==", "!=", "<", ">", "<=", ">=":
		return EQUALITY
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
	case "true", "false":
		return p.parseBooleanLiteral()
	case "null":
		return p.parseNullLiteral()
	case "undefined":
		return p.parseUndefinedLiteral()
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

func (p *Parser) parseBooleanLiteral() *BooleanLiteral {
	lit := &BooleanLiteral{Token: p.curToken}
	lit.Value = (p.curToken.Text == "true")
	return lit
}

func (p *Parser) parseNullLiteral() *NullLiteral {
	return &NullLiteral{Token: p.curToken}
}

func (p *Parser) parseUndefinedLiteral() *UndefinedLiteral {
	return &UndefinedLiteral{Token: p.curToken}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken() // Consume "("
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(")") {
		return nil
	}
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

func (p *Parser) curTokenIs(t string) bool {
	return p.curToken.Type.Label == t
}

/**
 * Error handling
 */

func (p *Parser) Errors() []string {
	return p.errors
}
