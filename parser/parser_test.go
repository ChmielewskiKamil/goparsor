package parser

import (
	"goparsor/ast"
	"goparsor/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
    let x = 5;
    let foo = 123;
    
    let bar = 0;
    `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram(...) returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"foo"},
		{"bar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.Statement, identName string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.tokenLiteral not 'let', got: %q", stmt.TokenLiteral())
		return false
	}

	// Assert the type of the stmt to be of type ast.LetStatement
	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", stmt)
		return false
	}
	if letStmt.Name.Value != identName {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", identName, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.TokenLiteral() != identName {
		t.Errorf("s.Name not '%s'. got=%s", identName, letStmt.Name)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
    return 5;

    return 0;
    return 99999999;

    `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram(...) returned nil")
	}

	stmts := program.Statements

	if len(stmts) != 3 {
		t.Errorf("Expected: 3 return statements, got: %d", len(stmts))
	}

	for _, stmt := range stmts {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected statement of type *ast.ReturnStatement, got: %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("Return statement expected to have: 'return' literal, got: %q", returnStmt.TokenLiteral())
		}
	}
}

////////////////////////////////////////////////////////////////////
//                             UTILS 			                  //
////////////////////////////////////////////////////////////////////

func checkParserErrors(t *testing.T, parser *Parser) {
	errors := parser.errors

	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser encountered: %d errors.", len(errors))

	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}
