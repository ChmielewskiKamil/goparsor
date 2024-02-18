package lexer

import "goparsor/token"

type Lexer struct {
	// Source code
	input string
	// Current position in input (points to current char)
	position int
	// Current reading position in input (after current char).
	// It always points to one character after the current char,
	// so its position +1
	readPosition int
	// Current char under examination
	// The byte type only supports ASCII. Pike uses the rune type to support
	// all unicode and UTF-8 characters. These unsupported chars can have
	// a size of multiple bytes and would require special handling.
	ch byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()

	return tok
}

/*~*~*~*~*~*~*~*~*~*~*~*~* Helper Functions ~*~*~*~*~*~*~*~*~*~*~*~*~*/

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// The zero is a null in ASCII, which for us means EOF or
		// that we didn't read anything yet. In this case its the former.
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tkType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tkType,
		Literal: string(ch),
	}
}
