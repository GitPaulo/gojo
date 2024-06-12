package parser

import (
	"fmt"
	"gojo/lexer"
	"strings"
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
	String() string
}

type Expression interface {
	Node
	expressionNode()
	String() string
}

type Declaration interface {
	Node
	declarationNode()
	String() string
}

type Node interface {
	TokenLiteral() string
	String() string
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out strings.Builder
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return "Program(" + out.String() + ")"
}

// VariableDeclaration represents a variable declaration.
type VariableDeclaration struct {
	Token lexer.GojoToken
	Name  *Identifier
	Value Expression
}

func (vd *VariableDeclaration) statementNode()       {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Text }
func (vd *VariableDeclaration) String() string {
	return fmt.Sprintf("VariableDeclaration(%s = %s)", vd.Name.String(), vd.Value.String())
}

// Identifier represents a variable name.
type Identifier struct {
	Token lexer.GojoToken
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Text }
func (i *Identifier) String() string {
	return i.Value
}

// IntegerLiteral represents an integer.
type IntegerLiteral struct {
	Token lexer.GojoToken
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Text }
func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("%d", il.Value)
}

// StringLiteral represents a string.
type StringLiteral struct {
	Token lexer.GojoToken
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Text }
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", sl.Value)
}

// BooleanLiteral represents a boolean.
type BooleanLiteral struct {
	Token lexer.GojoToken
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Text }
func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("%t", bl.Value)
}

// NullLiteral represents a null value.
type NullLiteral struct {
	Token lexer.GojoToken
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return nl.Token.Text }
func (nl *NullLiteral) String() string {
	return "null"
}

// UndefinedLiteral represents an undefined value.
type UndefinedLiteral struct {
	Token lexer.GojoToken
}

func (ul *UndefinedLiteral) expressionNode()      {}
func (ul *UndefinedLiteral) TokenLiteral() string { return ul.Token.Text }
func (ul *UndefinedLiteral) String() string {
	return "undefined"
}

// BinaryExpression represents a binary operation.
type BinaryExpression struct {
	Token    lexer.GojoToken
	Left     Expression
	Operator string
	Right    Expression
}

func (be *BinaryExpression) expressionNode()      {}
func (be *BinaryExpression) TokenLiteral() string { return be.Token.Text }
func (be *BinaryExpression) String() string {
	return fmt.Sprintf("BinaryExpression(%s %s %s)", be.Left.String(), be.Operator, be.Right.String())
}

// CallExpression represents a function call.
type CallExpression struct {
	Token     lexer.GojoToken
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Text }
func (ce *CallExpression) String() string {
	var args []string
	for _, arg := range ce.Arguments {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("CallExpression(%s(%s))", ce.Function.String(), strings.Join(args, ", "))
}

// BlockStatement represents a block of statements.
type BlockStatement struct {
	Token      lexer.GojoToken
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Text }
func (bs *BlockStatement) String() string {
	var out strings.Builder
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return "{" + out.String() + "}"
}

// FunctionDeclaration represents a function declaration.
type FunctionDeclaration struct {
	Token      lexer.GojoToken
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fd *FunctionDeclaration) statementNode()       {}
func (fd *FunctionDeclaration) TokenLiteral() string { return fd.Token.Text }
func (fd *FunctionDeclaration) String() string {
	var params []string
	for _, param := range fd.Parameters {
		params = append(params, param.String())
	}
	return fmt.Sprintf("FunctionDeclaration(%s(%s) %s)", fd.Name.String(), strings.Join(params, ", "), fd.Body.String())
}

// IfStatement represents an if-else statement.
type IfStatement struct {
	Token       lexer.GojoToken
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Text }
func (is *IfStatement) String() string {
	out := fmt.Sprintf("IfStatement(%s %s", is.Condition.String(), is.Consequence.String())
	if is.Alternative != nil {
		out += " else " + is.Alternative.String()
	}
	return out + ")"
}
