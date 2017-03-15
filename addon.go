package curscraft

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Debaru/goutils/fileutil"
	"github.com/mholt/archiver"
)

type Addon struct {
	Name           string
	Id             string
	PackageVersion string
	Upgradable     bool
	NewestVersion  string
}

func (a *Addon) getTocInfo(path string) {
	r, _ := os.Open(path)
	s := bufio.NewScanner(r)
	for s.Scan() {
		t := s.Text()
		// Package Version
		if strings.Contains(t, "X-Curse-Packaged-Version") {
			v := strings.Split(t, ":")
			a.PackageVersion = strings.TrimSpace(v[1])
		}

		// Package Version #2 (If no "X-Curse-Packaged-Version")
		if len(a.PackageVersion) == 0 && strings.Contains(t, "Version") {
			v := strings.Split(t, ":")
			a.PackageVersion = strings.TrimSpace(v[1])
		}

		// Project Name
		if strings.Contains(t, "X-Curse-Project-Name") {
			v := strings.Split(t, ":")
			a.Name = strings.TrimSpace(v[1])
		}

		// No project Name, use Title
		if len(a.Name) == 0 && strings.Contains(t, "Title") {
			v := strings.Split(t, ":")
			a.Name = strings.TrimSpace(v[1])
		}

		// Project ID
		if strings.Contains(t, "X-Curse-Project-ID") {
			v := strings.Split(t, ":")
			a.Id = strings.TrimSpace(v[1])
		}

		// Project ID #2 (If no "X-Curse-Project-ID")
		if len(a.Id) == 0 {
			r := strings.NewReplacer(" ", "-")
			a.Id = strings.ToLower(r.Replace(a.Name))
		}
	}
}

func (a *Addon) humanUpgradable() string {
	var h string

	h = fmt.Sprint("- Up to date")
	if a.Upgradable == true {
		h = fmt.Sprintf("- A new version is available (version %v)", a.NewestVersion)
	}

	return h
}

func getListAddons(path string) (addons []Addon, err error) {
	var urlAddon string

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !strings.Contains(file.Name(), "Blizzard_") {
			addon := new(Addon)
			addon.getTocInfo(path + "/" + file.Name() + "/" + file.Name() + ".toc")
			if len(addon.PackageVersion) > 0 {
				// Attribute Upgradable
				urlAddon = fmt.Sprintf("https://mods.curse.com/addons/wow/%v", addon.Id)
				tag := GetURL(urlAddon, "li.class=newest-file")
				if len(tag) > 0 {
					v, _ := tag[0].GetAttr("text")
					aVersion := strings.Split(v, ":")
					version := strings.TrimSpace(aVersion[1])
					if version != addon.PackageVersion {
						addon.Upgradable = true
						addon.NewestVersion = version
					}
				}

				addons = append(addons, *addon)
			}
		}
	}

	return addons, err
}

func DisplayAddonsList(path string) {
	addons, err := getListAddons(path)
	if err == nil {
		for _, addon := range addons {
			fmt.Println(addon.Name, addon.PackageVersion, addon.humanUpgradable())
		}
	} else {
		fmt.Println(err)
	}
}

func DownloadAddon(path string) {
	var urlDownload, urlAddon string

	addons, _ := getListAddons(path)
	for _, addon := range addons {
		urlAddon = fmt.Sprintf("https://mods.curse.com/addons/wow/%v", addon.Id)
		tag := GetURL(urlAddon, "li.class=newest-file")
		if len(tag) > 0 {
			v, _ := tag[0].GetAttr("text")
			aVersion := strings.Split(v, ":")
			version := strings.TrimSpace(aVersion[1])

			if version != addon.PackageVersion {
				urlDownload = fmt.Sprintf("https://mods.curse.com/addons/wow/%v/download", addon.Id)
				ta := GetURL(urlDownload, "a.class=download-link")
				test, _ := ta[0].GetAttr("data-href")
				f, err := fileutil.Download(test, path)
				defer os.Remove(f.Name())
				if err == nil {
					archiver.Zip.Open(f.Name(), path)
				}
			}
		}
	}
}
