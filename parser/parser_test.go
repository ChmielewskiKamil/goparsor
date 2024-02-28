package parser

import (
	"fmt"
	"goparsor/ast"
	"goparsor/lexer"
	"testing"
)

func TestParsingLetStatements(t *testing.T) {
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

func TestParsingReturnStatements(t *testing.T) {
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

func TestParsingIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected: 1 statement, got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected: *ast.ExpressionStatement, got: %T",
			program.Statements[0])
	}

	testIdentifier(t, stmt.Expression, "foobar")
}

func testIdentifier(t *testing.T, expr ast.Expression, value string) bool {
	ident, ok := expr.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected: *ast.Identifier, got: %T", expr)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("Expected token literal: '%s', got: %q",
			value, ident.TokenLiteral())
		return false
	}

	if ident.Value != value {
		t.Errorf("Expected ident value: '%s', got: %q", value, ident.Value)
		return false
	}

	return true
}

func TestParsingIntegerLiterals(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected: 1 statement, got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected: *ast.ExpressionStatement, got: %T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected: *ast.IntegerLiteral, got: %T",
			stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("Expected literal value: 5, got: %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("Expected token literal: '5', got: %s", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-125", "-", 125},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected: 1 statement, got: %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected: *ast.ExpressionStatement, got: %T", program.Statements[0])
		}

		prefixExpr, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Expected: *ast.PrefixExpression, got: %T", stmt.Expression)
		}

		if prefixExpr.Operator != tt.operator {
			t.Fatalf("Expected: %s operator, got: %s", tt.operator, prefixExpr.Operator)
		}

		if !testIntegerLiteral(t, prefixExpr.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected: 1 statement, got: %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected: *ast.ExpressionStatement, got: %T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Expected: *ast.InfixExpression, got: %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, expr.Left, tt.leftValue) {
			return
		}

		if expr.Operator != tt.operator {
			t.Errorf("Expected: %s operator, got: %s", tt.operator, expr.Operator)
		}

		if !testIntegerLiteral(t, expr.Right, tt.rightValue) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLiteral, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected: *ast.IntegerLiteral, got: %T", il)
		return false
	}

	if intLiteral.Value != value {
		t.Errorf("Expected int literal value: %d, got: %d", value, intLiteral.Value)
		return false
	}

	if intLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("Expected int token literal: %d, got %s", value, intLiteral.TokenLiteral())
		return false
	}
	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("Expected: %s, got: %s", tt.expected, actual)
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
