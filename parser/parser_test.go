package parser

import (
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
