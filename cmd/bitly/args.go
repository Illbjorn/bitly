package main

import (
	"fmt"
	"os"

	"github.com/illbjorn/basicli"
	"github.com/illbjorn/bitly/internal/safe"
	"github.com/illbjorn/echo"
)

type Args struct {
	Show       bool     `basicli:"show,s"`
	Hex        bool     `basicli:"hex,h"`
	Binary     bool     `basicli:"binary,b"`
	Debug      bool     `basicli:"debug,d"`
	Dec        bool     `basicli:"dec"`
	Positional []string `basicli:"-"`
	Help       Help     `basicli:"help"`
}

type Help struct{}

func (self Help) Exec() {
	os.Stderr.Write(safe.Stob(fmt.Sprintf(`
Welcome to %[1]sBitly%[2]s, your friendly neighborhood expression evaluation
companion!

USAGE
  To enter the REPL, just call Bitly!

FLAGS
  --show,   -s  Controls output of each individual step along the order of
                operations.
  --hex,    -h  Controls output of base-16 for calculation results.
  --binary, -b  Controls output of base-2 for calculation results.
  --dec         Controls output of base-10 for calculation results.
  --debug,  -d  Controls output of additional diagnostic data.
  --prompt, -p  String to prefix REPL lines with.
`, ansiCyan, ansiReset)))
	os.Exit(0)
}

func parseArgs() (args Args) {
	if err := basicli.Unmarshal(&args); err != nil {
		echo.Errorf("Failed to parse command line inputs: %s.", err)
		os.Exit(1)
	}
	// Dispatch to handle any subcommands(help)
	_ = basicli.Dispatch(&args)

	return args
}

const (
	ansiBlack  = "\03390m"
	ansiRed    = "\03391m"
	ansiGreen  = "\03392m"
	ansiYellow = "\03393m"
	ansiBlue   = "\03394m"
	ansiPurple = "\03395m"
	ansiCyan   = "\03396m"
	ansiWhite  = "\03397m"
	ansiReset  = "\033[0m"
)
