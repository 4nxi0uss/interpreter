package ast

import (
	"monkey/token"
	"testing"
)

func TestString(test *testing.T) {
	program := &Program{
		Statemens: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		test.Errorf("program.String() wrong. got=%q", program.String())
	}

}
