package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func testEvaluate(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Evaluate(program)
}

func TestEvaluateIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, o object.Object, expected int64) bool {
	result, ok := o.(*object.Integer)
	if !ok {
		t.Errorf("object expected : object.Integer, but was actual : %T", o)
		return false
	}
	if result.Value != expected {
		t.Errorf("result.Value expected : %d, but was actual : %d", expected, result.Value)
		return false
	}
	return true
}

func TestEvaluateBooleanObject(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, o object.Object, expected bool) bool {
	result, ok := o.(*object.Boolean)

	if !ok {
		t.Errorf("object expected : object.Boolean, but was actual : %T", o)
		return false
	}

	if result.Value != expected {
		t.Errorf("result.Value expected : %t, but was actual : %t", expected, result.Value)
		return false
	}
	return true
}
