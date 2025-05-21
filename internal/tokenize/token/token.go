package token

import "github.com/illbjorn/bitly/internal/ux"

func New(source []byte) Token {
	return Token{
		Meta: &ux.Meta{
			Source: source,
		},
	}
}

type Token struct {
	Kind Kind
	Meta *ux.Meta
}

func (self *Token) String() string {
	if self == nil {
		return ""
	}
	if self.Meta == nil {
		return ""
	}
	return self.Meta.Value()
}
