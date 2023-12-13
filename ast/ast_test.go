package ast

import (
	"monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.ID, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.ID, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	s := program.String()
	if s != "let myVar = anotherVar;" {
		t.Errorf("expected : [let myVar = anotherVar;], but was actual : [%s]", s)
	}
}
