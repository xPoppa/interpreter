package evaluator

import (
	"fmt"

	"github.com/xPoppa/interpreter/ast"
	"github.com/xPoppa/interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatements(node)
	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}
	}

	return nil
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}

	return result
}

func evalInfixExpression(operator string,
	left, right object.Object,
) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer)
	rightVal := right.(*object.Integer)

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal.Value + rightVal.Value}
	case "-":
		return &object.Integer{Value: leftVal.Value - rightVal.Value}
	case "*":
		return &object.Integer{Value: leftVal.Value * rightVal.Value}
	case "/":
		return &object.Integer{Value: leftVal.Value / rightVal.Value}
	case "<":
		return nativeBoolToBooleanObject(leftVal.Value < rightVal.Value)
	case ">":
		return nativeBoolToBooleanObject(leftVal.Value > rightVal.Value)
	case "==":
		return nativeBoolToBooleanObject(leftVal.Value == rightVal.Value)
	case "!=":
		return nativeBoolToBooleanObject(leftVal.Value != rightVal.Value)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, obj object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(obj)
	case "-":
		return evalMinusPrefixOperatorExpression(obj)
	default:
		return newError("unknown operator; %s %s", operator, obj.Type())
	}
}

func evalBangOperator(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(obj object.Object) object.Object {
	if obj.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", obj.Type())
	}

	value := obj.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
