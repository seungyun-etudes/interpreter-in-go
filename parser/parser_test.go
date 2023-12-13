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

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]
		if !testLetStatement(t, statement, tt.expectedIdentifier) {
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
