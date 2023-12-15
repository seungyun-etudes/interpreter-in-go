package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkProgramLength(t, 3, program)

	expected := []struct {
		identifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, e := range expected {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, e.identifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, expected string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral expected=%q, but was %q", "let", s.TokenLiteral())
		return false
	}
	letStatement, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s type expected=*ast.LetStatement, but was %T", s)
		return false
	}

	if letStatement.Name.Value != expected {
		t.Errorf("letStatement.Name.Value expected=*%s, but was %s", expected, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != expected {
		t.Errorf("letStatement.Name.TokenLiteral expected=*%s, but was %s", expected, letStatement.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	checkProgramLength(t, 3, program)

	for _, statement := range program.Statements {
		testReturnStatement(t, statement)
	}
}

func testReturnStatement(t *testing.T, statement ast.Statement) {
	returnStatement, ok := statement.(*ast.ReturnStatement)

	if !ok {
		t.Errorf("statement expected : *ast.ReturnStatement. but was actual : %T", statement)
	}

	if returnStatement.TokenLiteral() != "return" {
		t.Errorf("statement.TokenLiteral() expected : return, but was actual : %s", returnStatement.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkProgramLength(t, 1, program)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement expected : *ast.ExpressionStatement. but was actual : %T", program.Statements[0])
	}

	id, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Errorf("expression expected : *ast.Identifier. but was actual : %T", statement.Expression)
	}

	if id.Value != "foobar" {
		t.Errorf("id.Value expected : %s, but was actual : %s", "foobar", id.Value)
	}

	if id.TokenLiteral() != "foobar" {
		t.Errorf("id.TokenLiteral() expected : %s, but was actual : %s", "foobar", id.Value)
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := "5"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkProgramLength(t, 1, program)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement expected : *ast.ExpressionStatement. but was actual : %T", program.Statements[0])
	}

	numberLiteral, ok := statement.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Errorf("expression expected : *ast.NumberLiteral. but was actual : %T", statement.Expression)
	}

	if numberLiteral.Value != 5 {
		t.Errorf("numberLiteral.Value expected : %d, but was actual : %d", 5, numberLiteral.Value)
	}

	if numberLiteral.TokenLiteral() != "5" {
		t.Errorf("numberLiteral.TokenLiteral() expected : %s, but was actual : %s", "5", numberLiteral.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		checkProgramLength(t, 1, program)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("statement expected : *ast.ExpressionStatement. but was actual : %T", program.Statements[0])
		}

		expression := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("expression expected : *ast.PrefixExpression. but was actual : %T", statement.Expression)
		}

		if expression.Operator != tt.operator {
			t.Errorf("expression.Operator expected : %s, but was actual : %s", tt.operator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)
		checkProgramLength(t, 1, program)

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("statement expected : *ast.ExpressionStatement. but was actual : %T", program.Statements[0])
		}

		if !testInfixExpression(t, statement.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, e ast.Expression, left interface{}, operator string, right interface{}) bool {
	expression, ok := e.(*ast.InfixExpression)
	if !ok {
		t.Errorf("expression expected : *ast.InfixStatement. but was actual : %T", expression)
	}

	if !testLiteralExpression(t, expression.Left, left) {
		return false
	}

	if expression.Operator != operator {
		t.Fatalf("expression.Operator expected : %s, but was actual : %s", operator, expression.Operator)
	}

	if !testLiteralExpression(t, expression.Right, right) {
		return false
	}
	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 < 4 == 3 > 4", "((5 < 4) == (3 > 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"3 < 5 == true", "((3 < 5) == true)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if program.String() != tt.expected {
			t.Errorf("expected : %q, but was actual : %q", tt.expected, program.String())
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkProgramLength(t, 1, program)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("expression expected : ast.IfExpression, but was actual : %T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("Consequence len expected : 1, but was actual : %d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("expression.Consequence.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if expression.Alternative != nil {
		t.Errorf("expression.Alternative expected : nil, but was actual : %+v", expression.Alternative)
		return
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkProgramLength(t, 1, program)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("expression expected : ast.IfExpression, but was actual : %T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if len(expression.Consequence.Statements) != 1 {
		t.Errorf("Consequence len expected : 1, but was actual : %d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("expression.Consequence.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(expression.Alternative.Statements) != 1 {
		t.Errorf("Alternative len expected : 1, but was actual : %d", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("expression.Alternative.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `function(x, y) { x + y; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkProgramLength(t, 1, program)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("function expected ast.FunctionLiteral: ast.s, but was actual : %T", statement.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Errorf("function.Parameters len expected : 1, but was actual : %d", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Errorf("function.Body.Statements len expected : 1, but was actual : %d", len(function.Body.Statements))
	}

	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("function.Body.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)
	checkProgramLength(t, 1, program)

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements[0] expected : ast.ExpressionStatement, but was actual : %T", program.Statements[0])
	}

	call, ok := statement.Expression.(*ast.CallExpression)
	if !ok {
		t.Errorf("call expected : ast.CallExpression, but was actual : %T", statement.Expression)
	}

	if !testIdentifier(t, call.Function, "add") {
		return
	}

	if len(call.Arguments) != 3 {
		t.Errorf("call.Arguments len expected : 1, but was actual : %d", len(call.Arguments))
	}

	testLiteralExpression(t, call.Arguments[0], 1)
	testInfixExpression(t, call.Arguments[1], 2, "*", 3)
	testInfixExpression(t, call.Arguments[2], 4, "+", 5)
}

func testNumberLiteral(t *testing.T, expression ast.Expression, value int64) bool {
	numberLiteral, ok := expression.(*ast.NumberLiteral)

	if !ok {
		t.Errorf("expression expected : ast.NumberLiteral, but was actual : %T", expression)
		return false
	}

	if numberLiteral.Value != value {
		t.Errorf("numberLiteral.Value expected : %d, but was actual : %d", value, numberLiteral.Value)
		return false
	}

	if numberLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("numberLiteral.TokenLiteral() expected : %s, but was actual : %s", fmt.Sprintf("%d", value), numberLiteral.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, expression ast.Expression, value string) bool {
	id, ok := expression.(*ast.Identifier)

	if !ok {
		t.Errorf("expression expected : ast.Identifier, but was actual %T", expression)
		return false
	}

	if id.Value != value {
		t.Errorf("id.Value expected : %s, but was actual : %s", value, id.Value)
		return false
	}

	if id.TokenLiteral() != value {
		t.Errorf("id.TokenLiteral expected : %s, but was actual : %s", value, id.TokenLiteral())
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	b, ok := expression.(*ast.Boolean)
	if !ok {
		t.Errorf("exp expected : ast.Boolean, but was actual : %T", b)
		return false
	}

	if b.Value != value {
		t.Errorf("b.Value expected : %t, but was actual : %t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral() expected : %s, but was actual : %s", fmt.Sprintf("%t", value), b.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testNumberLiteral(t, expression, int64(v))
	case int64:
		return testNumberLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	}

	t.Errorf("type of expression not handled : %T", expression)
	return false
}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.Errors()) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(p.errors))
	for _, msg := range p.Errors() {
		t.Errorf("parser error : %q", msg)
	}
	t.FailNow()
}

func checkProgramLength(t *testing.T, expected int, program *ast.Program) {
	if len(program.Statements) != expected {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", expected, len(program.Statements))
	}
}
