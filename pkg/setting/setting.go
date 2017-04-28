// Package setting implements curscraft settings
package setting

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Debaru/curscraft/pkg/utils"
)

const (
	name = "settings"
)

// Settings is representation of settings file
type Settings struct {
	WoWFolder string
}

// SetPath adding path to configuration
func (s *Settings) SetPath(p string) error {
	if len(p) == 0 {
		return errors.New("No path name")
	}

	s.WoWFolder = filepath.Join(p, "Interface", "AddOns")
	f, err := os.Open(s.WoWFolder)
	if err != nil {
		return err
	}

	fInfo, _ := f.Stat()

	if !fInfo.IsDir() {
		return errors.New("Not a correct path")
	}

	return nil
}

// Load load settings file into Settings struct
func Load() *Settings {
	var s Settings
	utils.OpenFile(name, &s)
	return &s
}

// Save convert and save Settings struct on json format
func (s *Settings) Save() error {
	return utils.SaveFile(name, s)
}
