package tokenize

import (
	"bytes"
	"testing"

	"github.com/illbjorn/bitly/internal/assert"
)

func TestTokenizer(t *testing.T) {
	assert.SetT(t)
	tz := NewTokenizer([]byte("3 & 1 << 4"))

	t.Run("peek", func(t *testing.T) {
		// Wildly out of bounds
		assert.Equal(tz.Peek(-342), EOF)
		assert.Equal(tz.Peek(342), EOF)

		// First 3 bytes
		assert.Equal(tz.Peek(0), EOF)
		assert.Equal(tz.Peek(1), '3')
		assert.Equal(tz.Peek(2), ' ')
		assert.Equal(tz.Peek(3), '&')

		// Advance and try again
		tz.Adv()
		assert.Equal(tz.Peek(0), '3')

		// Advance and try again
		tz.Adv()
		assert.Equal(tz.Peek(0), ' ')
	})
	tz.Pos = -1

	t.Run("range", func(t *testing.T) {
		// Wildly out of range
		assert.True(tz.Range(-4, 327) == nil)
		assert.True(tz.Range(1, 327) == nil)

		// Low > High
		assert.True(tz.Range(4, 1) == nil)

		// In range
		assert.True(bytes.Equal(tz.Range(0, 3), []byte("3 &")))
		assert.True(bytes.Equal(tz.Range(1, 4), []byte(" & ")))
	})
	tz.Pos = -1

	t.Run("more", func(t *testing.T) {
		tz.Pos = -2
		assert.False(tz.More())

		tz.Adv()
		assert.True(tz.More())

		tz.Adv()
		assert.True(tz.More())

		for range len(tz.Input) - 1 {
			tz.Adv()
		}
		assert.False(tz.More())
	})
}
