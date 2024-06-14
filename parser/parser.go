package parser

import (
	"fmt"
	"gojo/config"
	"gojo/lexer"
	"strconv"
)

// Precedence Levels
const (
	LOWEST      = iota
	ASSIGN      // =
	CONDITIONAL // ?:
	LOGICAL_OR  // ||
	LOGICAL_AND // &&
	BITWISE_OR  // |
	BITWISE_XOR // ^
	BITWISE_AND // &
	EQUALS      // ==, !=
	COMPARISON  // <, >, <=, >=
	SHIFT       // <<, >>
	SUM         // +, -
	PRODUCT     // *, /, %
	PREFIX      // -X, !X
	CALL        // myFunction(X)
	MEMBER      // obj.property
	INDEX       // array[index]
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
	}

	p.peekToken = p.l.NextToken()

	if config.LoadConfig().Verbose {
		fmt.Println("Peek Token:", p.peekToken)
	}
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
	case "while":
		return p.parseWhileStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(";") {
		p.nextToken()
	}

	return stmt
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
	if leftExp == nil {
		return nil
	}

	if config.LoadConfig().Verbose {
		fmt.Println("╔══ parseExpression()")
		fmt.Println("Left Expression:", leftExp)
		fmt.Println("Precedence:", precedence)
		fmt.Println("Peek precedence:", p.peekPrecedence())
		fmt.Println("╚══ [entering loop = " + strconv.FormatBool(!p.peekTokenIs(";") && precedence < p.
			peekPrecedence()) + "]")
	}

	for !p.peekTokenIs(";") && precedence < p.peekPrecedence() {
		switch p.peekToken.Type.Label {
		case "(":
			p.nextToken()
			leftExp = p.parseCallExpression(leftExp)
		case ".":
			p.nextToken()
			leftExp = p.parseMemberAccessExpression(leftExp)
		case "=":
			p.nextToken()
			leftExp = p.parseAssignmentExpression(leftExp)
		default:
			if infixPrecedence := p.peekPrecedence(); precedence < infixPrecedence {
				p.nextToken()
				leftExp = p.parseInfixExpression(leftExp)
			} else {
				return leftExp
			}
		}
	}

	return leftExp
}

func (p *Parser) parseMemberAccessExpression(object Expression) Expression {
	exp := &MemberAccessExpression{Token: p.curToken, Object: object}

	if !p.expectPeek("identifier") {
		return nil
	}

	exp.Property = &Identifier{Token: p.curToken, Value: p.curToken.Text}
	return exp
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

func getPrecedence(tok lexer.GojoToken) int {
	switch tok.Type.Label {
	case "=":
		return ASSIGN
	case "||":
		return LOGICAL_OR
	case "&&":
		return LOGICAL_AND
	case "|":
		return BITWISE_OR
	case "^":
		return BITWISE_XOR
	case "&":
		return BITWISE_AND
	case "==", "!=":
		return EQUALS
	case "<", ">", "<=", ">=":
		return COMPARISON
	case "<<", ">>":
		return SHIFT
	case "+", "-":
		return SUM
	case "*", "/", "%":
		return PRODUCT
	case "(":
		return CALL
	case ".":
		return MEMBER
	case "[":
		return INDEX
	case "!":
		return PREFIX
	default:
		return LOWEST
	}
}

func (p *Parser) curPrecedence() int {
	return getPrecedence(p.curToken)
}

func (p *Parser) peekPrecedence() int {
	return getPrecedence(p.peekToken)
}

func (p *Parser) parseAtomicExpression() Expression {
	switch p.curToken.Type.Label {
	case "string":
		return p.parseStringLiteral()
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
	case "!":
		return p.parsePrefixExpression()
	case "while":
		return p.parseWhileStatement()
	default:
		return nil
	}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Text,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseAssignmentExpression(left Expression) Expression {
	exp := &AssignmentExpression{Token: p.curToken, Name: left.(*Identifier)}

	p.nextToken() // Move past '='
	exp.Value = p.parseExpression(LOWEST)

	return exp
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
	value := p.curToken.Type.Label == "true"
	return &BooleanLiteral{Token: p.curToken, Value: value}
}

func (p *Parser) parseNullLiteral() *NullLiteral {
	return &NullLiteral{Token: p.curToken}
}

func (p *Parser) parseUndefinedLiteral() *UndefinedLiteral {
	return &UndefinedLiteral{Token: p.curToken}
}

func (p *Parser) parseStringLiteral() *StringLiteral {
	return &StringLiteral{Token: p.curToken, Value: p.curToken.Text}
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken() // Consume "("
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(")") {
		return nil
	}
	return exp
}

func (p *Parser) parseWhileStatement() *WhileStatement {
	stmt := &WhileStatement{Token: p.curToken}

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

	stmt.Body = p.parseBlockStatement()

	return stmt
}

func (p *Parser) parseCallExpression(function Expression) Expression {
	exp := &CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(")")
	return exp
}

func (p *Parser) parseExpressionList(end string) []Expression {
	var list []Expression

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(",") {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.peekTokenIs(end) {
		p.errors = append(p.errors, fmt.Sprintf("expected next token to be %s, got %s instead", end, p.peekToken.Type.Label))
		return nil
	}

	p.nextToken() // consume the end token

	return list
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
