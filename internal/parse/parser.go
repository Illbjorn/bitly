package parse

import (
	"slices"

	"github.com/illbjorn/bitly/internal/tokenize/token"
	"github.com/illbjorn/bitly/internal/ux"
)

type Parser struct {
	Tokens []token.Token
	Pos    int
}

func (self *Parser) Peek(i int) *token.Token {
	i += self.Pos
	if i < 0 || i >= len(self.Tokens) {
		return nil
	}
	return &self.Tokens[i]
}

func (self *Parser) Adv() {
	self.Pos++
}

func (self *Parser) Want(kinds ...token.Kind) (*token.Token, error) {
	next := self.Peek(1)
	if next == nil {
		return nil, ErrEOF
	}

	if slices.Contains(kinds, next.Kind) {
		self.Adv()
		return next, nil
	}

	return nil, ux.
		NewError("unexpected token [%s]", next.Kind).
		WithSrc(next.Meta.Source).
		WithAnnotation(int(next.Meta.Sx), int(next.Meta.Ex))
}

func (self *Parser) Maybe(kinds ...token.Kind) (*token.Token, bool) {
	next := self.Peek(1)
	if next == nil {
		return nil, false
	}

	if slices.Contains(kinds, next.Kind) {
		self.Adv()
		return next, true
	}

	return nil, false
}

func (self *Parser) Snapshot() int {
	return self.Pos
}

func (self *Parser) Revert(i int) {
	self.Pos = i
}
