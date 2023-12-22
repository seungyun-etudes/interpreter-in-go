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
	return Evaluate(program, object.NewEnvironment())
}

func TestEvaluateIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == true", false},
		{"(1 < 2) == false", false},
		{"(1 > 2) == false", true},
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

func TestEvaluateBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvaluateIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, o object.Object) bool {
	if o != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", o, o)
		return false
	}
	return true
}

func TestEvaluateReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5 + true;", "type mismatch : INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch : INTEGER + BOOLEAN"},
		{"-true", "unknown operator : -BOOLEAN"},
		{"true + false", "unknown operator : BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator : BOOLEAN + BOOLEAN"},
		{"if(10 > 1) { true + false; }", "unknown operator : BOOLEAN + BOOLEAN"},
		{"if(10 > 1) { if(10 > 1) { return true + false; } return 1;}", "unknown operator : BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found : foobar"},
		{`"Hello" - "World"`, "unknown operator : STRING - STRING"},
	}

	for _, tt := range tests {
		evaluated := testEvaluate(tt.input)

		errorObject, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("evaluated expected : object.Error, but was actual : %T(%+v)", evaluated, evaluated)
			continue
		}

		if errorObject.Message != tt.expected {
			t.Errorf("errorObject.Message expected : %s, but was actual : %s", tt.expected, errorObject.Message)
		}
	}
}

func TestEvaluateLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a ;", 5},
		{"let a = 5 * 5; a ;", 25},
		{"let a = 5; let b = a ; b", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEvaluate(tt.input), tt.expected)
	}
}

func TestEvaluateFunctionObject(t *testing.T) {
	input := "function(x) { x + 2; }"

	evaluated := testEvaluate(input)
	function, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("evaluated expected : object.Function, but was actual %T", evaluated)
	}

	if len(function.Parameters) != 1 {
		t.Fatalf("parameter len expected : 1, but was actual %d", len(function.Parameters))
	}

	if function.Parameters[0].String() != "x" {
		t.Fatalf("parameter expected : x, but was actual %s", function.Parameters[0].String())
	}

	expectedBody := "(x + 2)"

	if function.Body.String() != expectedBody {
		t.Fatalf("body expected : %s, but was actual %s", expectedBody, function.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = function(x) { x; } identity(5);", 5},
		{"let identity = function(x) { return x; } identity(5);", 5},
		{"let double = function(x) { x * 2; } double(5);", 10},
		{"let add = function(x, y) { x + y; } add(5, 5);", 10},
		{"let add = function(x, y) { x + y; } add(5 + 5, add(5, 5));", 20},
		{"function(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEvaluate(tt.input), tt.expected)
	}
}

func TestEvaluateClosure(t *testing.T) {
	input := `
let newAdder = function(x) {
	function(y) {x + y;}
};

let addTwo = newAdder(2);
addTwo(5);
`

	testIntegerObject(t, testEvaluate(input), 7)
}

func TestEvaluateStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEvaluate(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Errorf("evaluated expected : object.String, but was actual : %T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("str.Value expected: Hello, World!, but was actual : %s", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!";`
	evaluated := testEvaluate(input)

	str, ok := evaluated.(*object.String)

	if !ok {
		t.Errorf("evaluated expected : object.String, but was actual : %T", evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("str.Value expected: Hello, World!, but was actual : %s", str.Value)
	}
}
