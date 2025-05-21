package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var (
	cfgDir, _    = os.UserConfigDir()
	bitlyCfgDir  = filepath.Join(cfgDir, "bitly")
	bitlyCfgPath = filepath.Join(bitlyCfgDir, "bitly.json")
)

func Load() error {
	f, err := os.Open(bitlyCfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(settings)
}

func Save() error {
	err := os.MkdirAll(bitlyCfgDir, 0o700)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(bitlyCfgPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(settings)
}
