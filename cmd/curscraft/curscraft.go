package main

import (
	"fmt"
	"os"

	"github.com/Debaru/curscraft"
)

const (
	curseUrl       = "https://mods.curse.com/addons/wow/"
	configDirName  = "curscraft"
	configFileName = "config.txt"
)

func list(s ...string) {
	curscraft.DisplayAddonsList(s[0])
}

func main() {
	var p map[string]interface{}
	p = make(map[string]interface{})

	// // Arguments List
	arguments := new(curscraft.Argts)
	arguments.Args = []curscraft.Arg{
		{Name: "set-path", Command: "-set-path", Shortcut: "-s", Help: "-set-path [PATH] : Path to your World Of Warcraft Addons repository.", Number: 1},
		{Name: "list", Command: "-list", Shortcut: "-l", Help: "-list : List of your World Of Warcraft Addons.", Number: 0},
		{Name: "update-all", Command: "-update-all", Shortcut: "-ua", Help: "-update-all : Update all your World Of Warcraft Addons.", Number: 0},
		{Name: "help", Command: "-help", Shortcut: "-h", Help: "-help : Display help.", Number: 0},
	}

	// Manage Arguments
	name, err := arguments.CheckArgument(os.Args)
	if err != nil {
		fmt.Println(err)
		arguments.Help()
		os.Exit(0)
	}
	// Configuration
	c := new(curscraft.Config)
	c.DirName = configDirName
	c.FileName = configFileName

	switch name {
	case "set-path":
		p["addon"] = os.Args[2]
		c.Properties = p
		err := c.Save()
		if err == nil {
			fmt.Println("Config File Update !")
		} else {
			fmt.Println(err)
		}

	case "list":
		c.Read()
		r := fmt.Sprint(c.Properties["addon"])
		curscraft.DisplayAddonsList(r)

	case "update-all":
		c.Read()
		r := fmt.Sprint(c.Properties["addon"])
		curscraft.DownloadAddon(r)
		fmt.Println("Addon are up to date.")

	default:
		arguments.Help()
	}

}
