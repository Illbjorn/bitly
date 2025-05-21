package parse

import "github.com/illbjorn/bitly/internal/ux"

var (
	ErrEOF        = ux.NewError("EOF")
	ErrNoMatch    = ux.NewError("no match found")
	ErrNotBinaryL = ux.NewError("not binary literal")
	ErrNotHexL    = ux.NewError("not hex literal")
)
