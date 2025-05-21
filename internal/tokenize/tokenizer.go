package tokenize

const EOF uint8 = '\x00'

func NewTokenizer(input []byte) *Tokenizer {
	return &Tokenizer{
		Input: input,
		Pos:   -1,
		Line:  1,
		Col:   1,
	}
}

type Tokenizer struct {
	Input     []byte
	Pos       int
	Line, Col uint32
}

func (self *Tokenizer) Adv() {
	if self.Peek(1) == '\n' {
		self.Line++
		self.Col = 1
	} else {
		self.Col++
	}
	self.Pos++
}

func (self *Tokenizer) More() bool {
	i := self.Pos + 1
	if i < 0 {
		return false
	}

	if int(i) >= len(self.Input) {
		return false
	}

	return true
}

func (self *Tokenizer) Peek(i int) byte {
	i += self.Pos

	// Out of bounds (low)
	if i < 0 {
		return EOF
	}

	// Out of bounds(high)
	if int(i) >= len(self.Input) {
		return EOF
	}

	return self.Input[i]
}

// Range basically handles the bounds check for substring slicing.
func (self *Tokenizer) Range(low, high int) []byte {
	// Out of bounds (low)
	if low < 0 {
		return nil
	}

	// Invalid range
	if low > high {
		return nil
	}

	// Out of bounds (high)
	if high > len(self.Input) {
		return nil
	}

	return self.Input[low:high]
}
