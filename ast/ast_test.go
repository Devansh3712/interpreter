package ast

import (
	"testing"

	"github.com/Devansh3712/interpreter/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "name"},
					Value: "name",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "tanu"},
					Value: "tanu",
				},
			},
		},
	}

	if program.String() != "let name = tanu;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
