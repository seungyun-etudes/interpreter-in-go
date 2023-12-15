package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Evaluate(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression)
	case *ast.PrefixExpression:
		return evaluatePrefixExpression(node.Operator, Evaluate(node.Right))
	case *ast.NumberLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.InfixExpression:
		return evaluateInfixExpression(node.Operator, Evaluate(node.Left), Evaluate(node.Right))
	case *ast.BlockStatement:
		return evaluateBlockStatement(node)
	case *ast.IfExpression:
		return evaluateIfExpression(node)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Evaluate(node.ReturnValue)}
	}
	return nil
}

func evaluateIfExpression(expression *ast.IfExpression) object.Object {
	condition := Evaluate(expression.Condition)

	if isTruthy(condition) {
		return Evaluate(expression.Consequence)
	} else if expression.Alternative != nil {
		return Evaluate(expression.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(o object.Object) bool {
	switch o {
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

func evaluateInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evaluateIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

func evaluateIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(right)
	case "-":
		return evaluateMinusPrefixOperatorExpression(right)
	default:
		return NULL
	}
}

func evaluateMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return NULL
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evaluateBangOperatorExpression(right object.Object) object.Object {
	switch right {
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

func nativeBoolToBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func evaluateProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Evaluate(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evaluateBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Evaluate(statement)

		if result != nil && result.Type() == object.RETURN_VALUE_OBJECT {
			return result
		}
	}

	return result
}
