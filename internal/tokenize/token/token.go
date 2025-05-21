package token

import "github.com/illbjorn/bitly/internal/safe"

func New(source []byte) Token {
	return Token{
		Meta: &Meta{
			Source: source,
		},
	}
}

type Token struct {
	Kind Kind
	Meta *Meta
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

type Meta struct {
	Start, Stop uint32
	Sx, Sy      uint32
	Ex, Ey      uint32
	Source      []byte
}

func (self Meta) Value() string {
	if self.Start < 0 { // BC
		return ""
	}
	if int(self.Stop) > len(self.Source) { // BC
		return ""
	}
	if self.Start > self.Stop { // BC
		return ""
	}
	slice := self.Source[self.Start:self.Stop]
	return safe.Btos(slice)
}
