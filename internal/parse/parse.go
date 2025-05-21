package parse

import (
	"github.com/illbjorn/bitly/internal/tokenize/token"
	"github.com/illbjorn/bitly/internal/ux"
)

func Parse(tokens []token.Token) (v any, err error) {
	p := &Parser{
		Tokens: tokens,
		Pos:    -1,
	}

	if v, err = ParseStatement(p); err == nil {
		return
	}

	return ParseExprAdd(p)
}

func ParseStatement(p *Parser) (stmt any, err error) {
	snap := p.Snapshot()

	// Bind
	if stmt, err = ParseBind(p); err == nil {
		return
	}

	// Set
	if stmt, err = ParseSet(p); err == nil {
		return
	}

	// Help
	if stmt, err = ParseHelp(p); err == nil {
		return
	}

	p.Revert(snap)

	return nil, ErrNoMatch
}

func ParseBind(p *Parser) (bind Bind, err error) {
	// bind
	//   : VAR ID EQ expr
	//   ;

	// VAR
	if _, err = p.Want(token.Var); err != nil {
		return
	}

	// ID
	if bind.ID, err = p.Want(token.ID); err != nil {
		return
	}

	// EQ
	if _, err = p.Want(token.EQ); err != nil {
		return
	}

	// expr
	if bind.Expr, err = ParseExpr(p); err != nil {
		return
	}

	return
}

func ParseSet(p *Parser) (set Set, err error) {
	snap := p.Snapshot()

	// SET
	if _, err = p.Want(token.Set); err != nil {
		p.Revert(snap)
		return
	}

	// ID
	if set.Name, err = p.Want(token.ID); err != nil {
		p.Revert(snap)
		return
	}

	// EQ
	if _, err = p.Want(token.EQ); err != nil {
		p.Revert(snap)
		return
	}

	// Basic
	if set.Value, err = ParseExprBasic(p); err != nil {
		p.Revert(snap)
		return
	}

	return
}

func ParseHelp(p *Parser) (help Help, err error) {
	snap := p.Snapshot()

	if _, err = p.Want(token.Help); err != nil {
		p.Revert(snap)
		return
	}

	return
}

func ParseExpr(p *Parser) (add Add, err error) {
	return ParseExprAdd(p)
}

func ParseExprAdd(p *Parser) (add Add, err error) {
	// Left
	if add.Left, err = ParseExprMult(p); err != nil {
		return add, err
	}

	for {
		var next AddM
		// Op
		if next.Op, err = p.Want(token.Plus, token.Minus); err != nil {
			return add, nil
		}

		// Right
		if next.Right, err = ParseExprMult(p); err != nil {
			return add, err
		}
		add.More = append(add.More, next)
	}
}

func ParseExprMult(p *Parser) (mult *Mult, err error) {
	mult = new(Mult)
	// Left
	if mult.Left, err = ParseExprPow(p); err != nil {
		return mult, err
	}

	for {
		var next MultM
		// Op
		if next.Op, err = p.Want(token.Asterisk, token.Slash, token.Percent); err != nil {
			return mult, nil
		}

		// Right
		if next.Right, err = ParseExprPow(p); err != nil {
			return mult, err
		}
		mult.More = append(mult.More, next)
	}
}

func ParseExprPow(p *Parser) (pow *Pow, err error) {
	pow = new(Pow)
	// Left
	if pow.Left, err = ParseExprShift(p); err != nil {
		return pow, err
	}

	for {
		var next PowM
		// Op
		if next.Op, err = p.Want(token.Caret); err != nil {
			return pow, nil
		}

		// Right
		if next.Right, err = ParseExprShift(p); err != nil {
			return pow, err
		}
		pow.More = append(pow.More, next)
	}
}

func ParseExprShift(p *Parser) (shift *Shift, err error) {
	shift = new(Shift)
	// Left
	if shift.Left, err = ParseExprLogic(p); err != nil {
		return shift, err
	}

	for {
		var next ShiftM
		// Op
		if next.Op, err = p.Want(token.LArrow2, token.RArrow2); err != nil {
			return shift, nil
		}

		// Right
		if next.Right, err = ParseExprLogic(p); err != nil {
			return shift, err
		}
		shift.More = append(shift.More, next)
	}
}

func ParseExprLogic(p *Parser) (logic *Logic, err error) {
	logic = new(Logic)
	// Left
	if logic.Left, err = ParseExprBasic(p); err != nil {
		return logic, err
	}

	for {
		var next LogicM
		// Op
		if next.Op, err = p.Want(token.Ampersand, token.Pipe); err != nil {
			return logic, nil
		}

		// Right
		if next.Right, err = ParseExprBasic(p); err != nil {
			return logic, err
		}
		logic.More = append(logic.More, next)
	}
}

func ParseExprBasic(p *Parser) (basic *Basic, err error) {
	basic = new(Basic)

	// SUB (-)(unary negate)
	basic.Negate, _ = p.Want(token.Minus)

	// Base-2
	if basic.BasicValue, err = p.Want(token.Base2); err == nil {
		return
	}

	// Base-10
	if basic.BasicValue, err = p.Want(token.Base10); err == nil {
		return
	}

	// Base-16
	if basic.BasicValue, err = p.Want(token.Base16); err == nil {
		return
	}

	// Identifier
	if basic.BasicValue, err = p.Want(token.ID); err == nil {
		return
	}

	// String
	if basic.BasicValue, err = p.Want(token.String); err == nil {
		return
	}

	// Group
	//
	// TODO

	// Init the error
	e := ux.NewError("expected operand").WithSrc(p.Tokens[0].Meta.Source)

	// Attempt to annotate
	next := p.Peek(1)
	if next = p.Peek(1); next != nil {
		e = e.WithAnnotation(int(next.Meta.Sx), int(next.Meta.Ex))
	} else if next = p.Peek(0); next != nil {
		e = e.WithAnnotation(int(next.Meta.Sx), int(next.Meta.Ex))
	}

	return nil, e
}
