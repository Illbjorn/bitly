package ux

import "github.com/illbjorn/bitly/internal/safe"

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
