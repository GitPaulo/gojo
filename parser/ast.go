package parser

import (
	"gojo/lexer"
)

// Program is the root node of every AST.
// It consists of a series of statements.
type Program struct {
	Statements []Statement
	Start      int
	End        int
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Declaration interface {
	Node
	declarationNode()
}

type Node interface {
	TokenLiteral() string
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// VariableDeclarationKeyword VariableDeclaration represents a variable declaration.
type VariableDeclarationKeyword string

const (
	KeywordVar VariableDeclarationKeyword = "var"
	KeywordLet VariableDeclarationKeyword = "let"
)

type VariableDeclaration struct {
	Token   lexer.GojoToken
	Name    *Identifier
	Value   Expression
	Keyword VariableDeclarationKeyword
}

func (ls *VariableDeclaration) statementNode()       {}
func (ls *VariableDeclaration) TokenLiteral() string { return ls.Token.Text }

// Identifier represents a variable name.
type Identifier struct {
	Token lexer.GojoToken
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Text }

// IntegerLiteral represents an integer.
type IntegerLiteral struct {
	Token lexer.GojoToken
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Text }

// StringLiteral represents a string.
type StringLiteral struct {
	Token lexer.GojoToken
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Text }

// BooleanLiteral represents a boolean.
type BooleanLiteral struct {
	Token lexer.GojoToken
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Text }

// NullLiteral represents a null value.
type NullLiteral struct {
	Token lexer.GojoToken
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return nl.Token.Text }

// UndefinedLiteral represents an undefined value.
type UndefinedLiteral struct {
	Token lexer.GojoToken
}

func (ul *UndefinedLiteral) expressionNode()      {}
func (ul *UndefinedLiteral) TokenLiteral() string { return ul.Token.Text }

// BinaryExpression represents a binary operation.
type BinaryExpression struct {
	Token    lexer.GojoToken
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode() {}

func (be *BinaryExpression) TokenLiteral() string {
	return be.Token.Text
}

// CallExpression represents a function call.
type CallExpression struct {
	Token     lexer.GojoToken
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}
