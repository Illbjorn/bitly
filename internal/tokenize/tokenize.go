package tokenize

import (
	"fmt"
	"strings"

	"github.com/illbjorn/bitly/internal/safe"
	"github.com/illbjorn/bitly/internal/tokenize/token"
	. "github.com/illbjorn/bitly/internal/tokenize/token"
	"github.com/illbjorn/bitly/internal/ux"
)

var ErrNoInput = fmt.Errorf("no input")

func Tokenize(input []byte) ([]Token, error) {
	if len(input) == 0 {
		return nil, ErrNoInput
	}

	var (
		tokens = make([]token.Token, 0, 32)
		tz     = NewTokenizer(input)
	)
	for {
		next, nextNext := tz.Peek(1), tz.Peek(2)
		if next == EOF {
			return tokens, nil
		}

		tk := token.New(tz.Input)
		markTokenStart(tz, &tk)

		switch {
		case isSpace(next):
			// Discard whitespace
			tz.Adv()

		case next == '0' && nextNext == 'b':
			// Base-2 Literal (order is important!)
			tz.Pos += 2        // Consume the prefix
			tk.Meta.Start += 2 // Update the start offset
			if err := base2(tz, &tk); err != nil {
				return nil, err
			}

		case next == '0' && nextNext == 'x':
			// Base-16 Literal (order is important!)
			tz.Pos += 2        // Consume the prefix
			tk.Meta.Start += 2 // Update the start offset
			if err := base16(tz, &tk); err != nil {
				return nil, err
			}

		case isBase10(next):
			// Base-10 Literal (order is important!)
			if err := base10(tz, &tk); err != nil {
				return nil, err
			}

		case isAlpha(next):
			// Keyword
			if err := keyword(tz, &tk); err != nil {
				return nil, err
			}

		case isSymbol(next):
			// Symbol
			if err := symbol(tz, &tk); err != nil {
				return nil, err
			}

		case isQuote(next):
			// Strings
			if err := stringL(tz, &tk); err != nil {
				return nil, err
			}

		default:
			return nil, ux.
				NewError("found unexpected symbol [%c]", next).
				WithSrc(input).
				WithAnnotation(int(tz.Col), int(tz.Col))
		}

		if tk.Kind > 0 {
			markTokenStop(tz, &tk)
			tokens = append(tokens, tk)
		}
	}
}

func markTokenStart(tz *Tokenizer, tk *Token) {
	tk.Meta.Start = uint32(tz.Pos + 1)
	tk.Meta.Sy = tz.Line
	tk.Meta.Sx = tz.Col
}

func markTokenStop(tz *Tokenizer, tk *Token) {
	if tk.Meta.Stop == 0 {
		tk.Meta.Stop = uint32(tz.Pos + 1)
	}
	tk.Meta.Sy = tz.Line
	tk.Meta.Sx = tz.Col
}

/*------------------------------------------------------------------------------
 * Keywords
 *----------------------------------------------------------------------------*/

func isAlpha(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

var keywords = map[string]token.Kind{
	"var":  token.Var,
	"set":  token.Set,
	"help": token.Help,
}

func keyword(tz *Tokenizer, tk *Token) (err error) {
	// Valid keyword characters are:
	// alpha : [A-Za-z]+ ;
	// num   : [0-9]+    ;
	// UNDERSCORE : '_'  ;
	//
	// But all keywords must start with an `alpha` code point (this is addressed
	// on the caller side)
	start := tz.Pos + 1
	for tz.More() {
		next := tz.Peek(1)
		if !(isAlpha(next) || isBase10(next) || next == '_') {
			break
		}
		tz.Adv()
	}
	stop := tz.Pos + 1
	// Classify the word
	word := tz.Input[start:stop]
	key := safe.Btos(word)
	kind, ok := keywords[key]
	if !ok {
		tk.Kind = ID
	} else {
		tk.Kind = kind
	}
	return
}

/*------------------------------------------------------------------------------
 * Numbers
 *----------------------------------------------------------------------------*/

func isBase2(b byte) bool {
	return b == '0' || b == '1'
}

func base2(tz *Tokenizer, tk *Token) (err error) {
	tk.Kind = token.Base2
	for tz.More() {
		if !isBase2(tz.Peek(1)) {
			break
		}
		tz.Adv()
	}
	return
}

func isBase10(b byte) bool {
	return b >= 0x30 && b <= 0x39
}

// TODO: Handle floats properly
func base10(tz *Tokenizer, tk *Token) (err error) {
	tk.Kind = Base10
	decimal := false
	for tz.More() {
		next := tz.Peek(1)
		switch {
		case next == '.' && !decimal:
			decimal = true

		case next == '.' && decimal:
			err = ux.
				NewError("found ill-formed decimal value").
				WithSrc(tz.Input).
				WithAnnotation(int(tz.Col), int(tz.Col))
			return

		case isBase10(next): // OK

		default:
			return
		}
		tz.Adv()
	}
	return
}

func isBase16(b byte) bool {
	return b >= 'a' && b <= 'f' || b >= 'A' && b <= 'F' || b >= '0' && b <= '9'
}

func base16(tz *Tokenizer, tk *Token) (err error) {
	tk.Kind = Base16
	for tz.More() {
		if !isBase16(tz.Peek(1)) {
			return
		}
		tz.Adv()
	}
	return
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\n' || b == '\r' || b == '\t'
}

func isSymbol(b byte) bool {
	const symbols = "<<>>-+/*%^&|=()"
	return strings.IndexByte(symbols, b) != -1
}

func symbol(tz *Tokenizer, tk *Token) error {
	next, nextNext := tz.Peek(1), tz.Peek(2)

	switch {
	//
	// 2-byte symbols
	//

	case next == '<' && nextNext == '<':
		tk.Kind = LArrow2
		tz.Adv()
		tz.Adv()

	case next == '>' && nextNext == '>':
		tk.Kind = RArrow2
		tz.Adv()
		tz.Adv()

		//
		// 1-byte symbols
		//

	case next == '&':
		tk.Kind = Ampersand
		tz.Adv()

	case next == '|':
		tk.Kind = Pipe
		tz.Adv()

	case next == '^':
		tk.Kind = Caret
		tz.Adv()

	case next == '%':
		tk.Kind = Percent
		tz.Adv()

	case next == '*':
		tk.Kind = Asterisk
		tz.Adv()

	case next == '/':
		tk.Kind = Slash
		tz.Adv()

	case next == '+':
		tk.Kind = Plus
		tz.Adv()

	case next == '-':
		tk.Kind = Minus
		tz.Adv()

	case next == '(':
		tk.Kind = ParenL
		tz.Adv()

	case next == ')':
		tk.Kind = ParenR
		tz.Adv()

	case next == '=':
		tk.Kind = EQ
		tz.Adv()

	default:
		return ux.
			NewError("found unexpected symbol [%c][%c]", next, nextNext).
			WithSrc(tz.Input).
			WithAnnotation(int(tz.Col), int(tz.Col))
	}

	return nil
}

func isQuote(b byte) bool {
	return b == '\'' || b == '"'
}

func stringL(tz *Tokenizer, tk *token.Token) error {
	tk.Kind = token.String

	term := tz.Peek(1)                 // ' | "
	tz.Adv()                           // ' | "
	tk.Meta.Start = uint32(tz.Pos + 1) // Adjust the start position to after the '
	for {
		if !tz.More() {
			return ux.
				NewError("reached EOF looking for matching [%c]", term).
				WithSrc(tz.Input).
				WithAnnotation(int(tk.Meta.Sx), int(tk.Meta.Ex))
		}
		if tz.Peek(1) == term {
			break
		}
		tz.Adv() // ' | "
	}
	tk.Meta.Stop = uint32(tz.Pos + 1)
	tz.Adv()

	return nil
}
