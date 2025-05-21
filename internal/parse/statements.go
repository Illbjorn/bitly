package parse

import "github.com/illbjorn/bitly/internal/tokenize/token"

type Statement = interface{ Bind | Set }

type Bind struct {
	ID   *token.Token
	Expr Expression
}

type Set struct {
	Name  *token.Token
	Value *Basic
}

type Help struct{}
