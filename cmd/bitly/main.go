package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/chzyer/readline"
	"github.com/illbjorn/bitly/internal/debug"
	"github.com/illbjorn/bitly/internal/parse"
	"github.com/illbjorn/bitly/internal/repl"
	"github.com/illbjorn/bitly/internal/settings"
	s "github.com/illbjorn/bitly/internal/settings"
	"github.com/illbjorn/bitly/internal/tokenize"
	"github.com/illbjorn/bitly/internal/tokenize/token"
	"github.com/illbjorn/echo"
)

func main() {
	// Parse command line inputs
	args := parseArgs()
	if args.Debug {
		// Set log flags
		// echo.SetFlags(
		// echo.WithColor,
		// echo.WithLevel,
		// echo.WithCallerFile,
		// echo.WithCallerLine,
		// )
		echo.SetLevel(echo.LevelDebug)
		echo.Debugf("Debug mode [ON].")
	}

	// Load Settings
	err := s.Load()
	if err != nil {
		echo.Errorf("Failed to load settings: %s.", err)
	}
	// Clobber from-file settings with CLI inputs
	MergeSettings(args)

	// Init the REPL
	r, err := repl.New(Prompt.Get())
	if err != nil {
		echo.Fatalf("Failed to init REPL: %s.", err)
	}
	defer r.Deinit()

	// Enter the REPL!
	for {
		// Wait for user input
		var line []byte
		if line, err = r.Read(); err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				err := s.Save()
				if err != nil {
					echo.Errorf("Failed to save settings: %s.", err)
				}
				return
			}
			echo.Fatalf("Failed to read input line: %s.", err)
		}

		// Tokenize
		tokens, err := Tokenize(line)
		if err != nil {
			if errors.Is(err, tokenize.ErrNoInput) {
				continue
			}
			echo.Errorf("Failed to tokenize input: %s", err)
			continue
		}

		// Parse
		expr, err := Parse(tokens)
		if err != nil {
			echo.Errorf("Failed to parse expression: %s", err)
			continue
		}

		// Evaluate
		err = Evaluate(expr, r)
		if err != nil {
			echo.Errorf("%s", err)
			continue
		}
	}
}

func Tokenize(input []byte) ([]token.Token, error) {
	if Debug.Get() {
		defer debug.ClockTo("Tokenize", os.Stderr)()
	}
	return tokenize.Tokenize(input)
}

func Parse(tokens []token.Token) (any, error) {
	if Debug.Get() {
		defer debug.ClockTo("Parse", os.Stderr)()
	}
	return parse.Parse(tokens)
}

func Evaluate(v any, r *repl.REPL) error {
	if Debug.Get() {
		defer debug.ClockTo("Evaluate", os.Stderr)()
	}

	switch v := v.(type) {
	case parse.Bind:
		i, err := Bind(v, r)
		if err != nil {
			return fmt.Errorf("Failed to bind [%s]: %s.", v.ID, err)
		}
		echo.Infof("[%s]=>[%d]", v.ID.Meta.Value(), i)
		return nil

	case parse.Set:
		err := ApplySetting(r, v)
		if err != nil {
			return fmt.Errorf("Failed to apply setting: %s.", err)
		}
		return nil

	case parse.Help:
		echo.Info(settings.List())
		return nil

	case parse.Expression:
		i, err := ExprAdd(v, r)
		if err != nil {
			return fmt.Errorf("Failed to evaluate expression: %s.", err)
		}

		if Bin.Get() {
			echo.Infof("Result[base-2]  => %#b", i)
		}

		if Dec.Get() {
			echo.Infof("Result[base-10] => %#d", i)
		}

		if Hex.Get() {
			echo.Infof("Result[base-16] => %#x", i)
		}
		return nil

	default:
		return fmt.Errorf("found unexpected evaluate struct kind [%T]", v)
	}
}
