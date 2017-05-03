package addon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"sync"

	"github.com/Debaru/curscraft/pkg/parser"
	"github.com/Debaru/curscraft/pkg/utils"
)

const (
	name = "data"
)

// Addons is a collection of Addon
type Addons []Addon

// Len for using sort package
func (slice Addons) Len() int {
	return len(slice)
}

// Less for using sort package
func (slice Addons) Less(i, j int) bool {
	return slice[i].Name < slice[j].Name
}

// Swap for using sort package
func (slice Addons) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Save list of addon in p file
func Save(addons *Addons) error {
	err := utils.SaveFile(name, addons)
	if err != nil {
		return err
	}
	return nil
}

// Load load addons from p file
// which is a reprensentation of all addons on json format
func Load() *Addons {
	var adds Addons
	utils.OpenFile(name, &adds)

	return &adds
}

// CheckAddons check all addons that are presents
// in wowFolder
func CheckAddons(addons *Addons) error {
	var wg sync.WaitGroup
	files, err := ioutil.ReadDir(WoWFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		wg.Add(1)
		go func(fileName string, addons *Addons) {
			defer wg.Done()
			var a Addon
			a.DirectoryName = fileName
			a.Real = true
			// Search if addon doesn't exist in addons
			i := searchByDirectoryName(addons, a.DirectoryName)
			if i == -1 {
				erro := a.readTocFile()
				if erro != nil {
					a.Real = false
				}

				// Check URL to Curse Web Site to set ID
				erro = a.setID()
				if erro != nil {
					a.Real = false
				}

				*addons = append(*addons, a)
			}
		}(file.Name(), addons)

	}
	wg.Wait()

	return nil
}

// List display Addon list
func List(addons *Addons) {
	sort.Sort(addons)
	for _, a := range *addons {
		if a.Real {
			fmt.Println(a.Name)
			fmt.Println("ID:", a.ID)
			fmt.Println("Author: ", a.Author)
			fmt.Println("Version: ", a.Version)
			if a.Upgradable {
				fmt.Println("New Version Available: ", a.NewVersion)
			}
			fmt.Print("URL: ", a.getCurseAddonURL(), "\n\n")
		}
	}
}

// Remove an addon
func Remove(addons *Addons, name string) error {
	// Searching Addon
	i := searchByID(addons, name)
	if i == -1 {
		return errors.New("Addon not found")
	}

	// Remove Addon Folder
	err := (*addons)[i].remove()
	if err != nil {
		return err
	}

	// Update list
	if i == 0 {
		(*addons) = append((*addons)[1:])
	} else {
		(*addons) = append((*addons)[:i], (*addons)[i+1:]...)
	}

	return nil
}

// CheckNewVersion checking if a new version is available
// return number of addon which have an update available
func CheckNewVersion(addons *Addons) int {
	var wg sync.WaitGroup
	c := 0
	for i, addon := range *addons {
		if addon.Real {
			wg.Add(1)
			go func(count *int, addon *Addon) {
				defer wg.Done()
				tag := parser.GetURL(addon.getCurseAddonURL(), "li.class=newest-file")
				if len(tag) > 0 {
					v, _ := tag[0].GetAttr("text")
					aVersion := strings.Split(v, ":")
					version := strings.TrimSpace(aVersion[1])
					if version != addon.Version {
						addon.Upgradable = true
						addon.NewVersion = version
						*count++
					}
				}
			}(&c, &(*addons)[i])
		}
	}
	wg.Wait()
	return c
}

// Upgrade upgrade all addon which are upgradable
// return number of addon which have been updated
func Upgrade(addons *Addons) (int, error) {
	c := 0
	for i := range *addons {
		if (*addons)[i].Upgradable && (*addons)[i].Real {
			err := (*addons)[i].download()
			if err != nil {
				return -1, err
			}
			c++
		}
	}

	return c, nil
}

// Download download addon with is url
// and check all new addons
func Download(addons *Addons, uri []string) error {
	var a Addon

	for i := 0; i < len(uri); i++ {
		url := uri[i]
		sp := strings.Split(url, "/")
		a.ID = sp[len(sp)-1]
		err := a.download()
		if err != nil {
			return err
		}
		fmt.Printf("Addon %s was downloaded.\n", a.ID)
	}

	// Check Addons
	err := CheckAddons(addons)
	if err != nil {
		return err
	}

	return err
}

func searchByDirectoryName(addons *Addons, dirName string) int {
	for i, a := range *addons {
		if a.DirectoryName == dirName {
			return i
		}
	}

	return -1
}

func searchByID(addons *Addons, ID string) int {
	for i, a := range *addons {
		if a.ID == ID {
			return i
		}
	}

	return -1
}
