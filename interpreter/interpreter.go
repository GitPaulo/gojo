package interpreter

import (
	"fmt"
	"gojo/parser"
)

type Interpreter struct {
	env map[string]interface{}
}

func New() *Interpreter {
	return &Interpreter{env: make(map[string]interface{})}
}

func (i *Interpreter) Interpret(program *parser.Program) {
	for _, stmt := range program.Statements {
		i.evalStatement(stmt)
	}
}

func (i *Interpreter) evalStatement(stmt parser.Statement) {
	switch stmt := stmt.(type) {
	case *parser.LetStatement:
		val := i.evalExpression(stmt.Value)
		i.env[stmt.Name.Value] = val
		fmt.Printf("%s = %v\n", stmt.Name.Value, val)
	}
}

func (i *Interpreter) evalExpression(expr parser.Expression) interface{} {
	switch expr := expr.(type) {
	case *parser.IntegerLiteral:
		return expr.Token.Text
	default:
		return nil
	}
}
