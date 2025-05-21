package repl

import (
	"github.com/chzyer/readline"
)

func New(prompt string) (r *REPL, err error) {
	r = new(REPL)

	// Init the bind map
	r.Binds = make(map[string]int64)

	// Init Readline
	if r.Instance, err = readline.New(prompt); err != nil {
		return nil, err
	}

	return
}

type REPL struct {
	Binds              map[string]int64 `json:"-"`
	*readline.Instance `json:"-"`
}

func (self *REPL) Read() ([]byte, error) {
	return self.Instance.ReadSlice()
}

func (self *REPL) Bind(name string, value int64) {
	self.Binds[name] = value
}

func (self *REPL) Lookup(name string) (i int64, ok bool) {
	i, ok = self.Binds[name]
	return
}

func (self *REPL) Deinit() error {
	return self.Instance.Close()
}
