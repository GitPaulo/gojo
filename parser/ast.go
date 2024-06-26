package parser

import (
	"fmt"
	"gojo/lexer"
	"strings"
)

/**
 * Interfaces
 */

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode() // Placeholder method to distinguish statements from expressions
}

type Expression interface {
	Node
	expressionNode() // Placeholder method to distinguish expressions from statements
}

/**
 * AST Nodes
 */

// Program represents the root of the AST.
type Program struct {
	Statements []Statement
	Start      int
	End        int
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
	Token      lexer.GojoToken
	Name       *Identifier
	Value      Expression
	IsConstant bool // Whether the variable is a "const"
}

func (vd *VariableDeclaration) statementNode()       {}
func (vd *VariableDeclaration) TokenLiteral() string { return vd.Token.Text }
func (vd *VariableDeclaration) String() string {
	return fmt.Sprintf("VariableDeclaration(%s %s = %s)", vd.Token.Type.Label, vd.Name.String(), vd.Value.String())
}

// AssignmentExpression represents a variable assignment.
type AssignmentExpression struct {
	Token lexer.GojoToken // The token (=)
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentExpression) expressionNode()      {}
func (ae *AssignmentExpression) TokenLiteral() string { return ae.Token.Text }
func (ae *AssignmentExpression) String() string {
	return fmt.Sprintf("AssignmentExpression(%s = %s)", ae.Name.String(), ae.Value.String())
}

// Identifier represents a variable name.
type Identifier struct {
	Token lexer.GojoToken
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Text }
func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier(%s)", i.Value)
}

// IntegerLiteral represents an integer.
type IntegerLiteral struct {
	Token lexer.GojoToken
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Text }
func (il *IntegerLiteral) String() string {
	return fmt.Sprintf("IntegerLiteral(%d)", il.Value)
}

// StringLiteral represents a string.
type StringLiteral struct {
	Token lexer.GojoToken
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Text }
func (sl *StringLiteral) String() string {
	return fmt.Sprintf("StringLiteral(\"%s\")", sl.Value)
}

// BooleanLiteral represents a boolean.
type BooleanLiteral struct {
	Token lexer.GojoToken
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Text }
func (bl *BooleanLiteral) String() string {
	return fmt.Sprintf("BooleanLiteral(%t)", bl.Value)
}

// NullLiteral represents a null value.
type NullLiteral struct {
	Token lexer.GojoToken
}

func (nl *NullLiteral) expressionNode()      {}
func (nl *NullLiteral) TokenLiteral() string { return nl.Token.Text }
func (nl *NullLiteral) String() string {
	return "NullLiteral(null)"
}

// UndefinedLiteral represents an undefined value.
type UndefinedLiteral struct {
	Token lexer.GojoToken
}

func (ul *UndefinedLiteral) expressionNode()      {}
func (ul *UndefinedLiteral) TokenLiteral() string { return ul.Token.Text }
func (ul *UndefinedLiteral) String() string {
	return "UndefinedLiteral(undefined)"
}

// ArrayLiteral represents an array.
type ArrayLiteral struct {
	Token    lexer.GojoToken
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Text }
func (al *ArrayLiteral) String() string {
	var elements []string
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	return fmt.Sprintf("ArrayLiteral(%s)", strings.Join(elements, ", "))
}

// ArrayAccessExpression represents an array access expression (e.g., arr[0]).
type ArrayAccessExpression struct {
	Token lexer.GojoToken
	Left  Expression // The array being accessed
	Index Expression // The index being accessed
}

func (aae *ArrayAccessExpression) expressionNode()      {}
func (aae *ArrayAccessExpression) TokenLiteral() string { return aae.Token.Text }
func (aae *ArrayAccessExpression) String() string {
	return fmt.Sprintf("ArrayAccessExpression(%s[%s])", aae.Left.String(), aae.Index.String())
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

// MemberAccessExpression represents a member access expression (e.g., obj.property).
type MemberAccessExpression struct {
	Token    lexer.GojoToken // The token (e.g., ".")
	Object   Expression      // The object being accessed
	Property *Identifier     // The property being accessed
}

func (mae *MemberAccessExpression) expressionNode()      {}
func (mae *MemberAccessExpression) TokenLiteral() string { return mae.Token.Text }
func (mae *MemberAccessExpression) String() string {
	return fmt.Sprintf("MemberAccessExpression(%s.%s)", mae.Object.String(), mae.Property.String())
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
	return fmt.Sprintf("CallExpression(%s(args=%s))", ce.Function.String(), strings.Join(args, ", "))
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

// WhileStatement represents a while loop.
type WhileStatement struct {
	Token     lexer.GojoToken
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Text }
func (ws *WhileStatement) String() string {
	return fmt.Sprintf("WhileStatement(%s, %s)", ws.Condition.String(), ws.Body.String())
}

// ExpressionStatement represents a statement consisting of a single expression.
type ExpressionStatement struct {
	Token      lexer.GojoToken // The first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Text }
func (es *ExpressionStatement) String() string {
	return fmt.Sprintf("ExpressionStatement(%s)", es.Expression.String())
}

// PrefixExpression represents a prefix operation (e.g., !true).
type PrefixExpression struct {
	Token    lexer.GojoToken // The prefix token, e.g., "!"
	Operator string          // The operator, e.g., "!"
	Right    Expression      // The expression to the right of the operator
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Text }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right.String())
}

// SwitchStatement represents a switch statement.
type SwitchStatement struct {
	Token       lexer.GojoToken
	Expression  Expression
	Cases       []*CaseClause
	DefaultCase *CaseClause
}

func (ss *SwitchStatement) statementNode()       {}
func (ss *SwitchStatement) TokenLiteral() string { return ss.Token.Text }
func (ss *SwitchStatement) String() string {
	var out strings.Builder
	out.WriteString("Switch (")
	out.WriteString(ss.Expression.String())
	out.WriteString(": ")
	for _, cc := range ss.Cases {
		out.WriteString(cc.String())
	}
	if ss.DefaultCase != nil {
		out.WriteString("Default: ")
		out.WriteString(ss.DefaultCase.String())
	}
	out.WriteString(")")
	return out.String()
}

// CaseClause represents a case clause in a switch statement.
type CaseClause struct {
	Token     lexer.GojoToken
	Condition Expression
	Body      *BlockStatement
}

func (cc *CaseClause) statementNode()       {}
func (cc *CaseClause) TokenLiteral() string { return cc.Token.Text }
func (cc *CaseClause) String() string {
	var out strings.Builder
	if cc.Condition != nil {
		out.WriteString("CaseClause(")
		out.WriteString(cc.Condition.String())
	} else {
		out.WriteString("DefaultCaseClause(")
	}
	out.WriteString("Body(")
	if cc.Body != nil {
		for _, stmt := range cc.Body.Statements {
			out.WriteString(stmt.String())
		}
	}
	out.WriteString(")) ")
	return out.String()
}

// BreakStatement represents a switch break statement.
type BreakStatement struct {
	Token lexer.GojoToken
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Text }
func (bs *BreakStatement) String() string {
	return "BreakStatement()"
}

// ReturnStatement represents a function/body return statement.
type ReturnStatement struct {
	Token lexer.GojoToken
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Text }
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("ReturnStatement(%s)", rs.Value.String())
}
