package parser

import (
	"fmt"
	"goparsor/ast"
	"goparsor/lexer"
	"goparsor/token"
	"strconv"
)

////////////////////////////////////////////////////////////////////
//                           CORE PARSER                          //
////////////////////////////////////////////////////////////////////

type Parser struct {
	l      *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.NextToken()
	p.NextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	return p
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

/*~*~*~*~*~*~*~*~*~*~*~*~* Pratt Parsing ~*~*~*~*~*~*~*~*~*~*~*~*~*/

// Operator precedence for Pratt Parsing
const (
	_           int = iota
	LOWEST          // Default, that we use for comparisons
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -A or !A
	CALL            // myFunction(A)
)

// Pratt Parsing functions
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

////////////////////////////////////////////////////////////////////
//                             Parsing  		                  //
////////////////////////////////////////////////////////////////////

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.NextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// We now that the current token is a statement
	stmt := &ast.LetStatement{Token: p.currToken}

	// The next token has to be identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// We already have the statement token (LET), now
	// we get its name e.g., let balance = 10
	// balance is the name
	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// @TODO: Implement expression handling
	// We skip the value of an expression for now.
	for !p.currTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currToken}

	// In the book there is NextToken() call here, but I don't think
	// it is necessary in this state.

	// @TODO: Skipped an expression part until a semicolon is found
	for !p.currTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// Semicolon is optional, that's why we either advance
	// or let it be, without throwing an error
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Check if there is an associated prefix parse function
	prefix := p.prefixParseFns[p.currToken.Type]
	if prefix == nil {
		return nil
	}

	// If there is such function, call it
	leftExpr := prefix()

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currToken}

	value, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse: %q as integer", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	literal.Value = value

	return literal
}

////////////////////////////////////////////////////////////////////
//                             UTILS 			                  //
////////////////////////////////////////////////////////////////////

func (p *Parser) currTokenIs(tkn token.TokenType) bool {
	return p.currToken.Type == tkn
}

func (p *Parser) peekTokenIs(tkn token.TokenType) bool {
	return p.peekToken.Type == tkn
}

func (p *Parser) expectPeek(tkn token.TokenType) bool {
	if p.peekTokenIs(tkn) {
		// The token type is what we expected it to be, so
		// just advance the lexer to the next token
		p.NextToken()
		return true
	} else {
		p.peekError(tkn)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(tkn token.TokenType) {
	msg := fmt.Sprintf("expected next token to be: %s, instead got: %s",
		tkn, p.peekToken.Type)

	p.errors = append(p.errors, msg)
}
