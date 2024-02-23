package parser

import (
	"goparsor/ast"
	"goparsor/lexer"
	"goparsor/token"
)

type Parser struct {
	l *lexer.Lexer

	currToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) NextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

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
	default:
		return nil
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
		return false
	}

}
