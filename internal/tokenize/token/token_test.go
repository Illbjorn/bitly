package token

import (
	"testing"

	"github.com/illbjorn/bitly/internal/assert"
)

func TestToken_Range(t *testing.T) {
	assert.SetT(t)

	tk := New([]byte("hello, world"))
	tk.Meta.Start = 0
	tk.Meta.Stop = 4
	assert.Equal(tk.Meta.Value(), "hell")
	tk.Meta.Start = 4
	tk.Meta.Stop = 6
	assert.Equal(tk.Meta.Value(), "o,")
}

func TestTokenKindString(t *testing.T) {
	assert.SetT(t)

	assert.Equal(None.String(), "<none>")
	assert.Equal(Base10.String(), "number")
	assert.Equal(Ampersand.String(), "&")
	assert.Equal(Pipe.String(), "|")
	assert.Equal(LArrow2.String(), "<<")
	assert.Equal(RArrow2.String(), ">>")
	assert.Equal(Asterisk.String(), "*")
	assert.Equal(Slash.String(), "/")
	assert.Equal(Plus.String(), "+")
	assert.Equal(Minus.String(), "-")
}
