package interpreter

import (
	"fmt"
	"gojo/config"
	"gojo/parser"
	"math"
	"strconv"
)

type Interpreter struct {
	Env map[string]interface{}
}

func New() *Interpreter {
	interpreter := &Interpreter{Env: make(map[string]interface{})}
	// Add built-in functions
	interpreter.Env["console"] = map[string]interface{}{
		"log": func(args ...interface{}) {
			fmt.Println(args...)
		},
	}
	interpreter.Env["Math"] = map[string]interface{}{
		"sqrt": func(x float64) float64 {
			return math.Sqrt(x)
		},
		"pow": func(x, y float64) float64 {
			return math.Pow(x, y)
		},
	}
	return interpreter
}

func (i *Interpreter) Interpret(program *parser.Program) {
	fmt.Println("╔═══ 🌸 Program Output:")
	for _, stmt := range program.Statements {
		i.evalStatement(stmt)
	}
	if config.LoadConfig().Verbose {
		fmt.Println("╔═══ 🌸 Program Environment:")
		maxKeyLength := 0
		for key := range i.Env {
			if len(key) > maxKeyLength {
				maxKeyLength = len(key)
			}
		}
		formatString := fmt.Sprintf("  %%-%ds: %%v\n", maxKeyLength)
		for key, val := range i.Env {
			fmt.Printf(formatString, key, val)
		}
	}
}

// InterpretREPL is used to interpret a single line of input in the REPL.
func (i *Interpreter) InterpretREPL(program *parser.Program) {
	for _, stmt := range program.Statements {
		i.evalStatement(stmt)
	}
}

func (i *Interpreter) evalStatement(stmt parser.Statement) {
	switch stmt := stmt.(type) {
	case *parser.VariableDeclaration:
		val := i.evalExpression(stmt.Value)
		i.Env[stmt.Name.Value] = val
	case *parser.IfStatement:
		i.evalIfStatement(stmt)
	case *parser.WhileStatement:
		for i.evalExpression(stmt.Condition).(bool) {
			i.evalBlockStatement(stmt.Body)
		}
	case *parser.ExpressionStatement:
		i.evalExpression(stmt.Expression)
	}
}

func (i *Interpreter) evalIfStatement(stmt *parser.IfStatement) {
	condition := i.evalExpression(stmt.Condition)
	if condition.(bool) {
		i.evalBlockStatement(stmt.Consequence)
	} else if stmt.Alternative != nil {
		switch alt := stmt.Alternative.Statements[0].(type) {
		case *parser.IfStatement:
			i.evalIfStatement(alt)
		default:
			i.evalBlockStatement(stmt.Alternative)
		}
	}
}

func (i *Interpreter) evalBlockStatement(block *parser.BlockStatement) {
	for _, stmt := range block.Statements {
		i.evalStatement(stmt)
	}
}

func (i *Interpreter) evalExpression(expr parser.Expression) interface{} {
	switch expr := expr.(type) {
	case *parser.StringLiteral:
		return expr.Value
	case *parser.BooleanLiteral:
		return expr.Value
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
	case *parser.AssignmentExpression:
		return i.evalAssignmentExpression(expr)
	case *parser.CallExpression:
		return i.evalCallExpression(expr)
	case *parser.MemberAccessExpression:
		object := i.evalExpression(expr.Object)
		if object == nil {
			fmt.Printf("Error: Object '%s' not found\n", expr.Object.String())
			return nil
		}
		if objName, ok := expr.Object.(*parser.Identifier); ok && objName.Value == "console" {
			// Handle 'console.log'
			if expr.Property.Value == "log" {
				return "console.log"
			} else {
				fmt.Printf("Error: Unsupported method '%s' for object 'console'\n", expr.Property.Value)
				return nil
			}
		}
		return object
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
			leftStr, leftOk := leftVal.(string)
			rightStr, rightOk := rightVal.(string)
			if leftOk && rightOk {
				return leftStr + rightStr
			}
			fmt.Println("Error: Invalid types for + operation")
			return nil
		case "-":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt - rightInt
			}
			fmt.Println("Error: Invalid types for - operation")
			return nil
		case "*":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt * rightInt
			}
			fmt.Println("Error: Invalid types for * operation")
			return nil
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
			fmt.Println("Error: Invalid types for / operation")
			return nil
		case "%":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt % rightInt
			}
			fmt.Println("Error: Invalid types for % operation")
			return nil
		case "<":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt < rightInt
			}
			leftStr, leftOk := leftVal.(string)
			rightStr, rightOk := rightVal.(string)
			if leftOk && rightOk {
				return leftStr < rightStr
			}
			fmt.Println("Error: Invalid types for < operation")
			return nil
		case ">":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt > rightInt
			}
			leftStr, leftOk := leftVal.(string)
			rightStr, rightOk := rightVal.(string)
			if leftOk && rightOk {
				return leftStr > rightStr
			}
			fmt.Println("Error: Invalid types for > operation")
			return nil
		case "<=":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt <= rightInt
			}
			leftStr, leftOk := leftVal.(string)
			rightStr, rightOk := rightVal.(string)
			if leftOk && rightOk {
				return leftStr <= rightStr
			}
			fmt.Println("Error: Invalid types for <= operation")
			return nil
		case ">=":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt >= rightInt
			}
			leftStr, leftOk := leftVal.(string)
			rightStr, rightOk := rightVal.(string)
			if leftOk && rightOk {
				return leftStr >= rightStr
			}
			fmt.Println("Error: Invalid types for >= operation")
			return nil
		case "==":
			return leftVal == rightVal
		case "!=":
			return leftVal != rightVal
		case "===":
			return leftVal == rightVal
		case "!==":
			return leftVal != rightVal
		case "&&":
			leftBool, leftOk := leftVal.(bool)
			rightBool, rightOk := rightVal.(bool)
			if leftOk && rightOk {
				return leftBool && rightBool
			}
			fmt.Println("Error: Invalid types for && operation")
			return nil
		case "||":
			leftBool, leftOk := leftVal.(bool)
			rightBool, rightOk := rightVal.(bool)
			if leftOk && rightOk {
				return leftBool || rightBool
			}
			fmt.Println("Error: Invalid types for || operation")
			return nil
		case "&":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt & rightInt
			}
			fmt.Println("Error: Invalid types for & operation")
			return nil
		case "|":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt | rightInt
			}
			fmt.Println("Error: Invalid types for | operation")
			return nil
		case "^":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt ^ rightInt
			}
			fmt.Println("Error: Invalid types for ^ operation")
			return nil
		case "<<":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt << rightInt
			}
			fmt.Println("Error: Invalid types for << operation")
			return nil
		case ">>":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return leftInt >> rightInt
			}
			fmt.Println("Error: Invalid types for >> operation")
			return nil
		case ">>>":
			leftInt, leftOk := leftVal.(int64)
			rightInt, rightOk := rightVal.(int64)
			if leftOk && rightOk {
				return int64(uint64(leftInt) >> uint64(rightInt))
			}
			fmt.Println("Error: Invalid types for >>> operation")
			return nil
		default:
			fmt.Printf("Error: Unsupported operator '%s'\n", expr.Operator)
			return nil
		}
	default:
		fmt.Println("Error: Unsupported expression type", expr)
	}
	return nil
}

func (i *Interpreter) evalCallExpression(expr *parser.CallExpression) interface{} {
	var functionName string
	switch fn := expr.Function.(type) {
	case *parser.Identifier:
		functionName = fn.Value
	case *parser.MemberAccessExpression:
		object := i.evalExpression(fn.Object)
		if object == nil {
			fmt.Printf("Error: Object '%s' not found\n", fn.Object.String())
			return nil
		}
		if objName, ok := fn.Object.(*parser.Identifier); ok {
			switch objName.Value {
			case "console":
				functionName = "console." + fn.Property.Value
			case "Math":
				functionName = "Math." + fn.Property.Value
			default:
				fmt.Printf("Error: Unsupported object '%s'\n", fn.Object.String())
				return nil
			}
		} else {
			fmt.Printf("Error: Unsupported object '%s'\n", fn.Object.String())
			return nil
		}
	default:
		fmt.Printf("Error: Unsupported function type '%s'\n", expr.Function.String())
		return nil
	}

	switch functionName {
	case "console.log":
		args := i.evalExpressions(expr.Arguments)
		if logFunc, ok := i.Env["console"].(map[string]interface{})["log"].(func(...interface{})); ok {
			logFunc(args...)
		} else {
			fmt.Println("Error: console.log function not found")
		}
		return nil
	case "Math.sqrt":
		if len(expr.Arguments) != 1 {
			fmt.Printf("Error: Math.sqrt expects 1 argument, got %d\n", len(expr.Arguments))
			return nil
		}
		arg := i.evalExpression(expr.Arguments[0])
		if num, ok := arg.(int64); ok {
			return math.Sqrt(float64(num))
		} else {
			fmt.Printf("Error: Math.sqrt expects a numeric argument\n")
			return nil
		}
	case "Math.pow":
		if len(expr.Arguments) != 2 {
			fmt.Printf("Error: Math.pow expects 2 arguments, got %d\n", len(expr.Arguments))
			return nil
		}
		base := i.evalExpression(expr.Arguments[0])
		exponent := i.evalExpression(expr.Arguments[1])
		if baseNum, ok1 := base.(int64); ok1 {
			if expNum, ok2 := exponent.(int64); ok2 {
				return math.Pow(float64(baseNum), float64(expNum))
			}
		}
		fmt.Printf("Error: Math.pow expects numeric arguments\n")
		return nil
	case "typeof":
		if len(expr.Arguments) != 1 {
			fmt.Printf("Error: typeof expects 1 argument, got %d\n", len(expr.Arguments))
			return nil
		}
		arg := i.evalExpression(expr.Arguments[0])
		return fmt.Sprintf("%T", arg)
	default:
		fmt.Printf("Error: Unsupported function call '%s'\n", functionName)
		return nil
	}
}

func (i *Interpreter) evalAssignmentExpression(expr *parser.AssignmentExpression) interface{} {
	_, ok := i.Env[expr.Name.Value]
	if !ok {
		fmt.Printf("Error (Line: %d): Variable '%s' not found\n", expr.Token.Line, expr.Name.Value)
		return nil
	}

	val := i.evalExpression(expr.Value)
	i.Env[expr.Name.Value] = val

	fmt.Printf("%s = %v (Line: %d)\n", expr.Name.Value, val, expr.Token.Line)
	return val
}

func (i *Interpreter) evalExpressions(exprs []parser.Expression) []interface{} {
	var result []interface{}
	for _, expr := range exprs {
		result = append(result, i.evalExpression(expr))
	}
	return result
}
