package ux

import (
	"bytes"
	"fmt"

	"github.com/illbjorn/bitly/internal/safe"
)

func NewError(msg string, values ...any) Error {
	if len(values) == 0 {
		return Error{
			Msg: msg,
		}
	}

	return Error{
		Msg: fmt.Sprintf(msg, values...),
	}
}

type Error struct {
	Msg               string
	Src               []byte
	BadStart, BadStop int
	Flags             Flags
}

var (
	buf = bytes.NewBuffer(make([]byte, 0, 1024))
)

const (
	caret  = '^'
	dashes = "-----------------------------------------------------------------" +
		"------------------------------------------------------------------------" +
		"------------------------------------------------------------------------" +
		"------------------------------------------------------------------------"
	spaces = "                                                                 " +
		"                                                                        " +
		"                                                                        " +
		"                                                                        "
	margin  = "  | "
	newline = '\n'
)

func (self Error) Error() string {
	buf.Reset()

	// Error message
	buf.WriteString(self.Msg)
	fmt.Fprintln(buf)

	// Source
	if self.Flags.WithSrc() {
		buf.WriteString(margin)
		buf.Write(self.Src)
		fmt.Fprintln(buf)
	}

	// Source Annotation
	if self.Flags.WithAnnotation() {
		buf.WriteString(margin)

		low, high := self.BadStart, self.BadStop

		// Indent
		i := min(
			max(0, low),
			len(spaces),
		) // BC
		buf.WriteString(spaces[:i])

		// Annotate
		diff := max(high-low, 0)
		switch {
		case diff < 2: // ^
			buf.WriteByte(caret)

		case diff == 2: // ^^
			buf.WriteByte(caret)
			buf.WriteByte(caret)

		case diff > 2: // ^-^
			buf.WriteByte(caret)
			buf.WriteString(dashes[:min(diff, len(dashes))])
			buf.WriteByte(caret)
		}
		buf.WriteByte('\n')
	}

	return safe.Btos(buf.Bytes())
}

func (self Error) WithSrc(src []byte) Error {
	self.Flags |= WithSrc
	self.Src = src
	return self
}

func (self Error) WithAnnotation(start, stop int) Error {
	self.Flags |= WithAnnotation
	self.BadStart, self.BadStop = start, stop
	return self
}

type Flags uint8

const (
	WithSrc Flags = 1 << iota
	WithAnnotation
)

func (self Flags) WithSrc() bool {
	return self&WithSrc == WithSrc
}

func (self Flags) WithAnnotation() bool {
	return self&WithAnnotation == WithAnnotation
}
