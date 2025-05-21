package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/illbjorn/bitly/internal/parse"
	"github.com/illbjorn/bitly/internal/repl"
	"github.com/illbjorn/bitly/internal/tokenize/token"
	"github.com/illbjorn/echo"
)

func ExprAdd(expr *parse.Add, r *repl.REPL) (int64, error) {
	left, err := ExprMult(expr.Left, r)
	if err != nil {
		return left, err
	}

	for _, next := range expr.More {
		right, err := ExprMult(next.Right, r)
		if err != nil {
			return left, err
		}

		left, err = Express(left, right, next.Op.Kind)
	}

	return left, err
}

func ExprMult(expr *parse.Mult, r *repl.REPL) (int64, error) {
	left, err := ExprPow(expr.Left, r)
	if err != nil {
		return left, err
	}

	for _, next := range expr.More {
		right, err := ExprPow(next.Right, r)
		if err != nil {
			return left, err
		}

		left, err = Express(left, right, next.Op.Kind)
	}

	return left, err
}

func ExprPow(expr *parse.Pow, r *repl.REPL) (int64, error) {
	left, err := ExprShift(expr.Left, r)
	if err != nil {
		return left, err
	}

	for _, next := range expr.More {
		right, err := ExprShift(next.Right, r)
		if err != nil {
			return left, err
		}

		left, err = Express(left, right, next.Op.Kind)
	}

	return left, err
}

func ExprShift(expr *parse.Shift, r *repl.REPL) (int64, error) {
	left, err := ExprLogic(expr.Left, r)
	if err != nil {
		return left, err
	}

	for _, next := range expr.More {
		right, err := ExprLogic(next.Right, r)
		if err != nil {
			return left, err
		}

		left, err = Express(left, right, next.Op.Kind)
	}

	return left, err
}

func ExprLogic(expr *parse.Logic, r *repl.REPL) (int64, error) {
	left, err := ExprBasic(expr.Left, r)
	if err != nil {
		return left, err
	}

	for _, next := range expr.More {
		right, err := ExprBasic(next.Right, r)
		if err != nil {
			return left, err
		}

		left, err = Express(left, right, next.Op.Kind)
	}

	return left, err
}

func ExprBasic(expr *parse.Basic, r *repl.REPL) (left int64, err error) {
	// Group
	if expr.Group != nil {
		left, err = ExprAdd(expr.Group, r)
		if expr.Negate != nil {
			return left * -1, err
		}
		return left, err
	}

	// Literal, variable
	switch expr.BasicValue.Kind {
	case token.Base2:
		left, err = strconv.ParseInt(expr.BasicValue.Meta.Value(), 2, 64)
		if err != nil {
			err = fmt.Errorf(
				"failed to parse base-2 value [%s]: %s",
				expr.BasicValue.Meta.Value(), err,
			)
			return
		}

	case token.Base10:
		left, err = strconv.ParseInt(expr.BasicValue.Meta.Value(), 10, 64)
		if err != nil {
			err = fmt.Errorf(
				"failed to parse base-10 value [%s]: %s",
				expr.BasicValue.Meta.Value(), err,
			)
			return
		}

	case token.Base16:
		left, err = strconv.ParseInt(expr.BasicValue.Meta.Value(), 16, 64)
		if err != nil {
			err = fmt.Errorf(
				"failed to parse base-16 value [%s]: %s",
				expr.BasicValue.Meta.Value(), err,
			)
			return
		}

	case token.ID:
		var ok bool
		left, ok = r.Lookup(expr.BasicValue.Meta.Value())
		if !ok {
			err = fmt.Errorf("found unknown bind [%s]", expr.BasicValue.Meta.Value())
			return
		}

	default:
		return 0, fmt.Errorf("found unexpected basic value kind [%s]", expr.BasicValue.Kind)
	}

	if expr.Negate != nil {
		return left * -1, nil
	}

	return left, nil
}

func Express(left, right int64, op token.Kind) (int64, error) {
	var result int64
	switch op {
	case token.None:
		return 0, fmt.Errorf(
			"Call to `Do` with no operator kind (%d)(%d).",
			left, right,
		)

	case token.Plus:
		result = left + right

	case token.Minus:
		result = left - right

	case token.Percent:
		result = left % right

	case token.Asterisk:
		result = left * right

	case token.Slash:
		result = left / right

	case token.Caret:
		result = int64(math.Pow(float64(left), float64(right)))

	case token.LArrow2:
		result = left << right

	case token.RArrow2:
		result = left >> right

	case token.Ampersand:
		result = left & right

	case token.Pipe:
		result = left | right
	}

	if Show.Get() {
		if Bin.Get() {
			echo.Infof("%#b %s %#b => %#b", left, op, right, result)
		} else if Hex.Get() {
			echo.Infof("%#x %s %#x => %#x", left, op, right, result)
		} else {
			echo.Infof("%d %s %d => %d", left, op, right, result)
		}
	}

	return result, nil
}
