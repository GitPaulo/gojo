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
		fmt.Printf("%s = %v\n", stmt.Name.Value, val)
	}
}

func (i *Interpreter) evalExpression(expr parser.Expression) interface{} {
	switch expr := expr.(type) {
	case *parser.IntegerLiteral:
		val, err := strconv.ParseInt(expr.Token.Text, 0, 64)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil
		}
		return val
	case *parser.BooleanLiteral:
		return expr.Value
	case *parser.Identifier:
		val, ok := i.Env[expr.Value]
		if !ok {
			fmt.Printf("Error: Variable '%s' not found\n", expr.Value)
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
					fmt.Println("Error: Division by zero")
					return nil
				}
			}
		case "&&":
			leftBool, leftOk := leftVal.(bool)
			rightBool, rightOk := rightVal.(bool)
			if leftOk && rightOk {
				return leftBool && rightBool
			}
		case "||":
			leftBool, leftOk := leftVal.(bool)
			rightBool, rightOk := rightVal.(bool)
			if leftOk && rightOk {
				return leftBool || rightBool
			}
		case "==":
			return leftVal == rightVal
		case "!=":
			return leftVal != rightVal
		default:
			fmt.Printf("Error: Unsupported operator '%s'\n", expr.Operator)
			return nil
		}
	default:
		fmt.Printf("Error: Unsupported expression type\n")
		return nil
	}
	return nil
}
