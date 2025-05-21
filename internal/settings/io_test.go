package settings

import (
	"path/filepath"
	"testing"

	"github.com/illbjorn/bitly/internal/assert"
)

func TestIO(t *testing.T) {
	assert.SetT(t)

	dir := t.TempDir()                                        //
	bitlyCfgPath = filepath.Join(dir, "cfg.json")             // Hijack setting path
	option := RegisterBool(Bool("an-option", true, true, "")) // Register an option
	assert.NoError(Save())                                    // Save settings
	option.Set(false)                                         // Mutate the setting
	assert.False(option.Get())                                // Verify set
	assert.NoError(Load())                                    // Load settings
	assert.True(option.Get())                                 // Confirm the load clobbered the manual assignment
	t.Log(bitlyCfgPath)
}
