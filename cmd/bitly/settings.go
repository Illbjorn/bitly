package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/illbjorn/bitly/internal/parse"
	"github.com/illbjorn/bitly/internal/repl"
	s "github.com/illbjorn/bitly/internal/settings"
)

// MergeSettings clobbers any `settings.Settings` (defined below) with values
// provided via CLI inputs.
func MergeSettings(args Args) {
	if args.Binary {
		Bin.Set(true)
	}
	if args.Hex {
		Hex.Set(true)
	}
	if args.Dec {
		Hex.Set(true)
	}
	if args.Show {
		Show.Set(true)
	}
	if args.Debug {
		Debug.Set(true)
	}
}

// ApplySetting accepts a REPL state and parsed `Set` statement, performing the
// appropriate setting assignment.
//
// The REPL state is passed in to allow mutation of the readline prompt.
func ApplySetting(r *repl.REPL, s *parse.Set) (err error) {
	key, value := s.Name.String(), s.Value.BasicValue.String()

	switch key {
	case "debug":
		var b bool
		b, err = strconv.ParseBool(value)
		if err == nil {
			Debug.Set(b)
			return nil
		}

	case "prompt":
		Prompt.Set(value)
		r.SetPrompt(value)
		return nil

	case "hex":
		var b bool
		b, err = strconv.ParseBool(value)
		if err == nil {
			Hex.Set(b)
			return nil
		}

	case "dec":
		var b bool
		b, err = strconv.ParseBool(value)
		if err == nil {
			Dec.Set(b)
			return nil
		}

	case "bin":
		var b bool
		b, err = strconv.ParseBool(value)
		if err == nil {
			Bin.Set(b)
			return nil
		}

	case "show":
		var b bool
		b, err = strconv.ParseBool(value)
		if err == nil {
			Show.Set(b)
			return nil
		}

	default:
		return fmt.Errorf("setting [%s] is invalid", key)
	}

	if err != nil {
		if errors.Is(err, strconv.ErrSyntax) {
			return fmt.Errorf(
				"[%s] is not a valid value for setting [%s]",
				value, key,
			)
		}
	}
	return
}

// Settings
var (
	Show = s.RegisterBool(s.Bool(
		"show",
		false,
		false,
		"Controls output of each individual step along the order of operations.",
	))
	Prompt = s.RegisterString(s.String(
		"prompt",
		"> ",
		"> ",
		"String to prefix REPL lines with.",
	))
	Debug = s.RegisterBool(s.Bool(
		"debug",
		false,
		false,
		"Controls output of additional diagnostic data.",
	))
	Hex = s.RegisterBool(s.Bool(
		"hex",
		false,
		false,
		"Controls output of base-16 for calculation results.",
	))
	Dec = s.RegisterBool(s.Bool(
		"dec",
		true,
		true,
		"Controls output of base-10 for calculation results.",
	))
	Bin = s.RegisterBool(s.Bool(
		"bin",
		false,
		false,
		"Controls output of base-2 for calculation results.",
	))
)
