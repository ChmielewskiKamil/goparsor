// The lexer implemented here is based
// on Rob Pike's talk on [Lexical Scanning in Go](https://www.youtube.com/watch?v=HxaD_trXwRE)
package lexer

import "fmt"

// The itemType identifies the type of lex items.
type itemType int

// The item represents a token returned from the scanner.
type item struct {
	typ itemType // type, such as itemNumber.
	val string   // value, such as "55.8".
}

// This funciton is known to printf which makes it good for debugging.
func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		// a quoted string with max of 10 characters, followed by "..."
		// e.g. input "Hello, World!", output "Hello, Wor"...
		return fmt.Sprintf("%.10q...", i.val)
	}

	return fmt.Sprintf("%q", i.val)
}

// lex type values.
const (
	itemError      itemType = iota // error occurred; value is text of error
	itemEOF                        // end of file
	itemDot                        // the cursor, spelled "."
	itemElse                       // the keyword "else"
	itemIf                         // the keyword "if"
	itemContract                   // the keyword "contract"
	itemText                       // plain text
	itemString                     // quoted string
	itemLeftBrace                  // the character "{"
	itemRightBrace                 // the character "}"
)

// The state represents where we are in the input and what we expect to see next.
// An action defines what we are going to do in that state given the input.
// After you execute the action, you will be in a new state.
// Combining the state and the action together results in a state function.
// The stateFn represents the state of the lexer as a function that returns the next state.
// It is a recursive definition.
type stateFn func(*lexer) stateFn

// The `run` function lexes the input by executing state functions
// until the state is nil.
func (l *lexer) run() {
	// @TODO: lexText is the first state function, usually
	// the start of the file. Change this to something meaningful.
	for state := lexText; state != nil; {
		state = state(l)
	}
	// The lexer is done, so we close the channel.
	// It tells the caller (probably the parser),
	// that no more tokens will be delivered.
	close(l.items)
}

// The lexer holds the state of the scanner.
type lexer struct {
	name  string    // Used only for error reports.
	input string    // The string being scanned.
	start int       // Start position of this item; in a big string, this is the start of the current token.
	pos   int       // Current position in the input.
	width int       // Width of last rune read from input.
	items chan item // Channel of scanned items.
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run()

	return l, l.items
}

// The `emit` function passes an item back to the client.
func (l *lexer) emit(t itemType) {
	// The value is a slice of the input.
	l.items <- item{t, l.input[l.start:l.pos]}
	// Move ahead in the input after sending it to the caller.
	l.start = l.pos
}
