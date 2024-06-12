package interpreter

import (
	"fmt"
	"gojo/parser"
	"strconv"
)

type Interpreter struct {
	Env map[string]interface{}
}

func New() *Interpreter {
	return &Interpreter{Env: make(map[string]interface{})}
}

func (i *Interpreter) Interpret(program *parser.Program) {
	fmt.Println("Interpreter Output:")
	fmt.Println("--------------------")
	for _, stmt := range program.Statements {
		i.evalStatement(stmt)
	}
}

func (i *Interpreter) evalStatement(stmt parser.Statement) {
	switch stmt := stmt.(type) {
	case *parser.VariableDeclaration:
		val := i.evalExpression(stmt.Value)
		i.Env[stmt.Name.Value] = val
		fmt.Printf("%s = %v (Line: %d)\n", stmt.Name.Value, val, stmt.Token.Line)
	}
}

func (i *Interpreter) evalExpression(expr parser.Expression) interface{} {
	switch expr := expr.(type) {
	case *parser.IntegerLiteral:
		val, err := strconv.ParseInt(expr.Token.Text, 0, 64)
		if err != nil {
			fmt.Printf("Error (Line: %d): %v\n", expr.Token.Line, err)
			return nil
		}
		return val
	case *parser.Identifier:
		val, ok := i.Env[expr.Value]
		if !ok {
			fmt.Printf("Error (Line: %d): Variable '%s' not found\n", expr.Token.Line, expr.Value)
			return nil
		}
		return val
	case *parser.BinaryExpression:
		leftVal := i.evalExpression(expr.Left)
		rightVal := i.evalExpression(expr.Right)
		switch expr.Operator {
		case "+":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt + rightInt
			}
		case "-":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt - rightInt
			}
		case "*":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt * rightInt
			}
		case "/":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				if rightInt != 0 {
					return leftInt / rightInt
				} else {
					fmt.Printf("Error (Line: %d): Division by zero\n", expr.Token.Line)
					return nil
				}
			}
		default:
			fmt.Printf("Error (Line: %d): Unsupported operator '%s'\n", expr.Token.Line, expr.Operator)
			return nil
		}
	default:
		fmt.Printf("Error: Unsupported expression type\n")
		return nil
	}
	return nil
}
