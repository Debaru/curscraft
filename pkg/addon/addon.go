package addon

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
}

func (a *Addon) setID() error {
	if len(a.ID) == 0 {
		r := strings.NewReplacer(" ", "-")
		a.ID = strings.ToLower(r.Replace(a.Name))
	}

	url := fmt.Sprintf(CURSE_URL, a.ID)
	r, err := http.Get(url)

	if err != nil {
		return err
	}

	// If code isn't 200, url must be with FileName
	if r.StatusCode != 200 {
		a.ID = strings.ToLower(a.DirectoryName)
		url = fmt.Sprintf(CURSE_URL, a.ID)
	}

	// Testing with new url
	r, err = http.Get(url)
	if err != nil {
		return err
	}

	// Addon no existing on Curse Website
	if r.StatusCode != 200 {
		return errors.New("No URL available")
	}

	// No Error Page
	t := parser.GetURL(url, "h2.text=Error")
	if len(t) > 0 {
		return errors.New("Addon URL unreachable")
	}

	return nil
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
