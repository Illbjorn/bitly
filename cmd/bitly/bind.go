package main

import (
	"github.com/illbjorn/bitly/internal/parse"
	"github.com/illbjorn/bitly/internal/repl"
)

func Bind(bind parse.Bind, r *repl.REPL) (i int64, err error) {
	i, err = ExprAdd(bind.Expr, r)
	r.Bind(bind.ID.Meta.Value(), i)
	return
}
