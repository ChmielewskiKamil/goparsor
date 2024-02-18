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

	l.skipWhitespace()

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
	default:
		if isLetter(l.ch) {
			// This literal is used to perform a lookup if this
			// particular identifier is a keyword or not.
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// Exit early because readIdentifier(...) advances l.position.
			// We don't want to do this again after switch statement.
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()

	return tok
}

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

func (l *Lexer) skipWhitespace() {
	// We don't care about whitespaces, newlines, carriage returns and tabs.
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\r' || l.ch == '\t' {
		// When encountered with any of these, just move forward.
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

/*~*~*~*~*~*~*~*~*~*~*~*~* Helper Functions ~*~*~*~*~*~*~*~*~*~*~*~*~*/

func newToken(tkType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tkType,
		Literal: string(ch),
	}
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}
