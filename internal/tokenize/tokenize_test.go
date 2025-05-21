package tokenize

import (
	"testing"

	"github.com/illbjorn/bitly/internal/assert"

	"github.com/illbjorn/bitly/internal/tokenize/token"
	. "github.com/illbjorn/bitly/internal/tokenize/token"
)

func TestTokenize(t *testing.T) {
	assert.SetT(t)
}

func TestIsNum(t *testing.T) {
	assert.SetT(t)

	// Wildly incorrect
	assert.False(isBase10('a'))

	for _, r := range "0123456789" {
		assert.True(isBase10(byte(r)))
	}
}

func TestNum(t *testing.T) {
	assert.SetT(t)

	// Number
	tz := NewTokenizer([]byte("123"))
	tk := token.New(tz.Input)
	err := base10(tz, &tk)
	assert.NoError(err)
	assert.Equal(tk.Kind, Base10)

	// Float with too many '.'s (error)
	tz = NewTokenizer([]byte("123.11.4"))
	tk = token.New(tz.Input)
	err = base10(tz, &tk)
	assert.Error(err)
}

func IsSpace(t *testing.T) {
	assert.SetT(t)

	// Wildly incorrect
	assert.False(isSpace('a'))

	assert.True(isSpace(' '))
	assert.True(isSpace('\n'))
	assert.True(isSpace('\r'))
	assert.True(isSpace('\t'))
}

func TestIsSymbol(t *testing.T) {
	assert.SetT(t)

	// Wildly incorrect
	assert.False(isSymbol('a'))

	assert.True(isSymbol('&'))
	assert.True(isSymbol('|'))
	assert.True(isSymbol('<'))
	assert.True(isSymbol('>'))
	assert.True(isSymbol('*'))
	assert.True(isSymbol('/'))
	assert.True(isSymbol('+'))
	assert.True(isSymbol('-'))
	assert.True(isSymbol('('))
	assert.True(isSymbol(')'))
}

func TestSymbol(t *testing.T) {
	assert.SetT(t)

	tz := NewTokenizer([]byte("<< >>&|*/+-"))
	tk := Token{}
	err := symbol(tz, &tk) // <<
	assert.Equal(tk.Kind, LArrow2)
	assert.NoError(err)

	tk = Token{}
	err = symbol(tz, &tk) // Whitespace (error)
	assert.Equal(tk.Kind, 0)
	assert.Error(err)
	tz.Adv()

	tk = Token{}
	err = symbol(tz, &tk) // >>
	assert.Equal(tk.Kind, RArrow2)
	assert.NoError(err)

	tk = Token{}
	err = symbol(tz, &tk) // &
	assert.Equal(tk.Kind, Ampersand)
	assert.NoError(err)

	tk = Token{}
	err = symbol(tz, &tk) // |
	assert.Equal(tk.Kind, Pipe)
	assert.NoError(err)

	tk = Token{}
	err = symbol(tz, &tk) // *
	assert.Equal(tk.Kind, Asterisk)
	assert.NoError(err)

	tk = Token{}
	err = symbol(tz, &tk) // /
	assert.Equal(tk.Kind, Slash)
	assert.NoError(err)

	tk = Token{}
	err = symbol(tz, &tk) // +
	assert.Equal(tk.Kind, Plus)
	assert.NoError(err)

	tk = Token{}
	err = symbol(tz, &tk) // -
	assert.Equal(tk.Kind, Minus)
	assert.NoError(err)
}
