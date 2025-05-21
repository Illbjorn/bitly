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
	Add        BinaryExpr[Mult]
	AddM       = BinaryExprMore[Mult]
	Mult       BinaryExpr[Pow]
	MultM      = BinaryExprMore[Pow]
	Pow        BinaryExpr[Shift]
	PowM       = BinaryExprMore[Shift]
	Shift      BinaryExpr[Logic]
	ShiftM     = BinaryExprMore[Logic]
	Logic      BinaryExpr[Basic]
	LogicM     = BinaryExprMore[Basic]

	Basic struct {
		Negate, BasicValue *token.Token
	}
)
