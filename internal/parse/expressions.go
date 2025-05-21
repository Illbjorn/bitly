package parse

import "github.com/illbjorn/bitly/internal/tokenize/token"

type BinaryExpr[T any] struct {
	Left *T
	More []BinaryExprMore[T]
}

type BinaryExprMore[T any] struct {
	Op    *token.Token
	Right *T
}

type (
	Expression = Add
	Add        BinaryExpr[Mult]        // mult
	AddM       = BinaryExprMore[Mult]  // ((+ | -) mult)*
	Mult       BinaryExpr[Pow]         // pow
	MultM      = BinaryExprMore[Pow]   // ((* | / | %) pow)*
	Pow        BinaryExpr[Shift]       // shift
	PowM       = BinaryExprMore[Shift] // (^ shift)*
	Shift      BinaryExpr[Logic]       // logic
	ShiftM     = BinaryExprMore[Logic] // ((<< | >>) logic)*
	Logic      BinaryExpr[Basic]       // basic
	LogicM     = BinaryExprMore[Basic] // ((& | |) basic)*

	Basic struct {
		Negate, BasicValue *token.Token
		Group              *Add
	}
)
