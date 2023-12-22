package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Evaluate(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evaluateProgram(node.Statements, environment)
	case *ast.ExpressionStatement:
		return Evaluate(node.Expression, environment)
	case *ast.PrefixExpression:
		right := Evaluate(node.Right, environment)
		if isError(right) {
			return right
		}
		return evaluatePrefixExpression(node.Operator, right)
	case *ast.NumberLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.InfixExpression:
		left := Evaluate(node.Left, environment)
		right := Evaluate(node.Right, environment)

		if isError(left) {
			return left
		}
		if isError(right) {
			return right
		}

		return evaluateInfixExpression(node.Operator, Evaluate(node.Left, environment), Evaluate(node.Right, environment))
	case *ast.BlockStatement:
		return evaluateBlockStatement(node, environment)
	case *ast.IfExpression:
		condition := Evaluate(node.Condition, environment)
		if isError(condition) {
			return condition
		}
		return evaluateIfExpression(node, environment)
	case *ast.ReturnStatement:
		value := Evaluate(node.ReturnValue, environment)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: Evaluate(node.ReturnValue, environment)}
	case *ast.LetStatement:
		value := Evaluate(node.Value, environment)
		if isError(value) {
			return value
		}
		environment.Set(node.Name.Value, value)
	case *ast.Identifier:
		return evaluateIdentifier(node, environment)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Environment: environment, Body: node.Body}
	case *ast.CallExpression:
		function := Evaluate(node.Function, environment)
		if isError(function) {
			return function
		}
		args := evaluateExpressions(node.Arguments, environment)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}

	return nil
}

func applyFunction(f object.Object, args []object.Object) object.Object {
	function, ok := f.(*object.Function)
	if !ok {
		return newError("not a function : %s", f.Type())
	}
	extendedEnvironment := extendFunctionEnvironment(function, args)
	evaluated := Evaluate(function.Body, extendedEnvironment)
	return unwrapReturnValue(evaluated)
}

func unwrapReturnValue(o object.Object) object.Object {
	if returnValue, ok := o.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return o
}

func extendFunctionEnvironment(function *object.Function, args []object.Object) *object.Environment {
	environment := object.NewEnclosedEnvironment(function.Environment)

	for i, param := range function.Parameters {
		environment.Set(param.Value, args[i])
	}

	return environment
}

func evaluateExpressions(expressions []ast.Expression, environment *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range expressions {
		evaluated := Evaluate(e, environment)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func isError(o object.Object) bool {
	if o != nil {
		return o.Type() == object.ERROR_OBJECT
	}
	return false
}

func evaluateIdentifier(node *ast.Identifier, environment *object.Environment) object.Object {
	value, ok := environment.Get(node.Value)

	if !ok {
		return newError("identifier not found : " + node.Value)
	}

	return value
}

func evaluateIfExpression(expression *ast.IfExpression, environment *object.Environment) object.Object {
	condition := Evaluate(expression.Condition, environment)

	if isTruthy(condition) {
		return Evaluate(expression.Consequence, environment)
	} else if expression.Alternative != nil {
		return Evaluate(expression.Alternative, environment)
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
	case left.Type() != right.Type():
		return newError("type mismatch : %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator : %s %s %s", left.Type(), operator, right.Type())
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
		return newError("unknown operator : %s %s %s", left.Type(), operator, right.Type())
	}
}

func evaluatePrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangOperatorExpression(right)
	case "-":
		return evaluateMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evaluateMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJECT {
		return newError("unknown operator : -%s", right.Type())
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

func evaluateProgram(statements []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Evaluate(statement, environment)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evaluateBlockStatement(block *ast.BlockStatement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Evaluate(statement, environment)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJECT || result.Type() == object.ERROR_OBJECT {
				return result
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
