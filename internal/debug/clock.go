package debug

import (
	"fmt"
	"io"
	"os"
	"time"
)

var stack = new(Stack)

func ClockTo(name string, w io.Writer) func() Result {
	stack.Frames = append(stack.Frames, Frame{
		Name:  name,
		Start: time.Now(),
		w:     w,
	})
	return Done
}

func Clock(name string) {
	stack.Frames = append(stack.Frames, Frame{
		Name:  name,
		Start: time.Now(),
	})
}

func Done() Result {
	if len(stack.Frames) == 0 {
		return Result{}
	}

	top := stack.Frames[len(stack.Frames)-1]
	r := Result{
		Name:     top.Name,
		Duration: time.Since(top.Start),
	}
	stack.Frames = stack.Frames[:len(stack.Frames)-1]

	if top.w != nil {
		r.WriteTo(top.w)
	}

	return r
}

type Stack struct {
	Frames []Frame
}

type Frame struct {
	Name  string
	Start time.Time
	w     io.Writer
}

type Result struct {
	Name     string
	Duration time.Duration
}

func (self Result) String() string {
	return fmt.Sprintf("  %-20s%-20s\n", self.Name, self.Duration)
}

func (self Result) WriteTo(w io.Writer) (n int64, err error) {
	v, err2 := os.Stderr.Write([]byte(self.String()))
	if err2 != nil {
		return n, err2
	}
	n += int64(v)

	return
}
