package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/Debaru/curscraft/pkg/addon"
	"github.com/Debaru/curscraft/pkg/setting"
	"github.com/Debaru/curscraft/pkg/utils"
)

var curscraftFolder string
var usage = `USAGE : curscraft <command> [argument]
DESCRIPTION
	-p, -path
		Path to your World of Warcraft Folder
	-l, -list
		List your World of Warcraft Addons
	-update
		Check if a new version is available
	-upgrade
		Upgrade addons which have a new version available
	-a, -add <addon URI...>
		Add addon(s)
	-r, -remove <addon ID>
		Remove addon

EXEMPLES
	-p, -path
		curscraft -path "/media/Blizzard/World of Warcraft"
	-a, -add
		curscraft -add https://mods.curse.com/addons/wow/recount
	-r, -remove
		curscraft -remove recount


`

func fatal(a ...interface{}) {
	fmt.Println(a)
	os.Exit(0)
}

func main() {
	var err error

	// Flag available
	var flagAvailable = sort.StringSlice{"-path", "-p", "-list", "-l", "-update", "-upgrade", "-r", "-remove", "-a", "-add"}
	flagAvailable.Sort()

	// Check Flag
	if len(os.Args) == 1 {
		fmt.Println(usage)
		os.Exit(0)
	}

	f := os.Args[1]
	i := flagAvailable.Search(f)
	if i == flagAvailable.Len() || flagAvailable[i] != f {
		fatal(usage)
	}

	// Load Setting File
	set := setting.Load()
	addon.WoWFolder = set.WoWFolder

	// Load Data (Addons) File
	adds := addon.Load()

	// Set Path
	if f == "-path" || f == "-p" {
		if len(os.Args) < 3 {
			fatal(usage)
		}

		err = set.SetPath(os.Args[2])
		utils.Fatal(err)
		addon.WoWFolder = set.WoWFolder
		err = set.Save()
		utils.Fatal(err)
		fmt.Println("World Of Warcraft folder is now setting.") // No error, display message

		// No addons, we're checking WoW Folder
		if adds.Len() == 0 {
			fmt.Println("Now, we're checking your addons. It can take a while...")
			err = addon.CheckAddons(adds)
			utils.Fatal(err)
			fmt.Printf("%d addon(s) found.\n", adds.Len())
		}

	}

	// Display List of Addons
	if f == "-list" || f == "-l" {
		addon.List(adds)
	}

	// Remove an Addon
	if f == "-remove" || f == "-r" {
		if len(os.Args) < 3 {
			fatal(usage)
		}

		err = addon.Remove(adds, os.Args[2])
		utils.Fatal(err)
		fmt.Println("Addon deleted.")

	}

	// Check new version
	if f == "-update" {
		fmt.Println("Checking for updates...")
		c := addon.CheckNewVersion(adds)
		fmt.Printf("%d addon(s) have an update. \n", c)
	}

	// Upgrade All Addon
	if f == "-upgrade" {
		fmt.Println("Now, we're upgrading your addons. It can take a while...")
		c, errU := addon.Upgrade(adds)
		utils.Fatal(errU)
		fmt.Printf("%d addon(s) have been updated. \n", c)

	}

	// Add Addon
	if f == "-add" || f == "-a" {
		if len(os.Args) < 3 {
			fatal(usage)
		}

		err = addon.Download(adds, os.Args[2:len(os.Args)])
		utils.Fatal(err)
	}

	// Save Addons
	err = addon.Save(adds)
	utils.Fatal(err)
}
