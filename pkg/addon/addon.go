package addon

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Debaru/curscraft/pkg/parser"
	"github.com/Debaru/goutils/fileutil"
	"github.com/mholt/archiver"
)

const (
	CURSE_URL    = "https://mods.curse.com/addons/wow/%v"
	CURSE_URL_DL = "https://mods.curse.com/addons/wow/%v/download"
)

// WoWFolder Path to WoW Folder
var WoWFolder string

// Addon reprsent an Addon Curse
type Addon struct {
	ID            string
	Name          string
	Version       string
	Author        string
	Upgradable    bool
	NewVersion    string
	DirectoryName string
	Real          bool
	dependencies  string
}

func (a *Addon) setID() error {
	var altID []string
	if len(a.ID) > 0 { // String array of alternative ID
		altID = append(altID, a.ID) // ID found in toc file
	}
	altID = append(altID, strings.ToLower(a.DirectoryName)) // Name of addon dir

	// Master Plan, become master-plan
	r := strings.NewReplacer(" ", "-")
	altID = append(altID, strings.ToLower(r.Replace(a.Name))) // Name of addon dir

	// BadBoy, become bad-boy
	str := []string{}
	for a, e := range []rune(a.Name) {
		b := unicode.IsUpper(e)

		if b && a > 0 {
			str = append(str, "-")
		}
		str = append(str, string(e))

	}

	title := strings.Join(str, "")
	altID = append(altID, strings.ToLower(title))

	for i := range altID {
		url := fmt.Sprintf(CURSE_URL, altID[i])
		r, err := http.Get(url)

		if err != nil {
			return err
		}

		if r.StatusCode == 200 {
			a.ID = altID[i]
			return nil
		}

	}

	return errors.New("No URL available")
}

func (a *Addon) setProperName() {
	s := strings.Split(a.Name, "|")
	a.Name = strings.TrimSpace(s[0])
}

func (a *Addon) getCurseAddonURL() string {
	return fmt.Sprintf(CURSE_URL, a.ID)
}

func (a *Addon) getCurseAddonDlURL() string {
	return fmt.Sprintf(CURSE_URL_DL, a.ID)
}

func (a *Addon) remove() error {
	p := filepath.Join(WoWFolder, a.DirectoryName)
	err := os.RemoveAll(p)
	if err != nil {
		return err
	}
	return nil
}

func (a *Addon) download() error {
	ta := parser.GetURL(a.getCurseAddonDlURL(), "a.class=download-link")
	v, err := ta[0].GetAttr("data-href")
	if err != nil {
		return err
	}

	f, err := fileutil.Download(v, WoWFolder)
	if err != nil {
		return err
	}

	defer os.Remove(f.Name())
	err = archiver.Zip.Open(f.Name(), WoWFolder)
	if err != nil {
		return err
	}

	// Update Addon info
	a.Upgradable = false
	a.Version = a.NewVersion

	return err
}
