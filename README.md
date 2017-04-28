# Curscraft

Curscraft is a program written in Go which allow to manage your addons for World of Warcraft.

## Installation

    git clone github.com/Debaru/curscraft
    go build curscraft/cmd/curscraft.go

## Usage
    ./curscraft <command> [argument]

 ## Commands
    	-p, -path
    		Path to your World of Warcraft Folder
    	-l, -list
    		List your World of Warcraft Addons
    	-update
    		Check if a new version is available
    	-upgrade
    		Upgrade addons which have a new version available
    	-a, -add <addon URL>
    		Add addon
    	-r, -remove <addon ID>
    		Remove addon


## Exemples
    	-p, -path
    		curscraft -path "/media/Blizzard/World of Warcraft"
    	-a, -add
    		curscraft -add https://mods.curse.com/addons/wow/deadly-boss-mods
    	-r, -remove
    		curscraft -remove recount

Configuration and Data files are stored here :
* Linux : ~/.config/curscraft/
