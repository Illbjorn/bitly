package token

import "github.com/illbjorn/echo"

type Kind uint8

const (
	None Kind = iota

	// "Primitives"
	Base10
	Base2
	Base16
	String

	// Arithmetic Operators
	Plus
	Minus
	Percent
	Asterisk
	Slash
	Caret
	LArrow2
	RArrow2
	Ampersand
	Pipe

	// Misc Symbols
	EQ
	ParenL
	ParenR
	Quote

	// Keywords
	Var
	Set
	Help

	// Identifier
	ID
)

// Guard for stringer support
var _ = [1]uint8{}[ID-22]

func (self Kind) String() string {
	switch self {
	case None:
		return "<none>"
	case Base10:
		return "<base-10>"
	case Base2:
		return "<base-2>"
	case Base16:
		return "<base-16>"
	case String:
		return "<string>"
	case Ampersand:
		return "&"
	case Pipe:
		return "|"
	case LArrow2:
		return "<<"
	case RArrow2:
		return ">>"
	case Caret:
		return "^"
	case Percent:
		return "%"
	case Asterisk:
		return "*"
	case Slash:
		return "/"
	case Plus:
		return "+"
	case Minus:
		return "-"
	case ParenL:
		return "("
	case ParenR:
		return ")"
	case Quote:
		return "\""
	case Var:
		return "var"
	case Set:
		return "set"
	case Help:
		return "help"
	case EQ:
		return "="
	case ID:
		return "<ID>"
	default:
		echo.Fatalf("Found unexpected token kind ['%d'].", self)
		panic("")
	}
}
