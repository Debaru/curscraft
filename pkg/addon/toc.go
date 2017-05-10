package addon

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (a *Addon) readTocFile() error {
	toc := fmt.Sprintf("%s.toc", a.DirectoryName)
	p := filepath.Join(WoWFolder, a.DirectoryName, toc)
	f, err := os.Open(p)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(f)
	for s.Scan() {
		t := s.Text()
		// Package Version
		if strings.Contains(t, "## X-Curse-Packaged-Version") {
			v := strings.Split(t, ":")
			a.Version = strings.TrimSpace(v[1])
		}

		// Package Version #2 (If no "X-Curse-Packaged-Version")
		if len(a.Version) == 0 && strings.Contains(t, "## Version") {
			v := strings.Split(t, ":")
			a.Version = strings.TrimSpace(v[1])
		}

		// Project Name
		if strings.Contains(t, "## X-Curse-Project-Name") {
			v := strings.Split(t, ":")
			a.Name = strings.TrimSpace(v[1])
		}

		// No project Name, use Title
		if len(a.Name) == 0 && strings.Contains(t, "## Title") {
			v := strings.Split(t, ":")
			a.Name = strings.TrimSpace(v[1])
		}

		// Project ID
		if strings.Contains(t, "## X-Curse-Project-ID") {
			v := strings.Split(t, ":")
			a.ID = strings.TrimSpace(v[1])
		}

		// Dependencies
		if strings.Contains(t, "## Dependencies") {
			v := strings.Split(t, ":")
			a.dependencies = strings.TrimSpace(v[1])
		}

	}

	return nil
}
