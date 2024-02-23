package ast

import (
	"goparsor/token"
	"testing"
)

// Make sure that the returned string matches the AST built by hand
// we are testing this piece of code:
// let monkey = chimpanzee
func TestASTString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "monkey"},
					Value: "monkey",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "chimpanzee"},
					Value: "chimpanzee",
				},
			},
		},
	}

	if program.String() != "let monkey = chimpanzee;" {
		t.Errorf("Program string is wrong, got: %q", program.String())
	}
}
