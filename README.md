# Curscraft

Curscraft is a program written in Go which allow to keep update your addons for World of Warcraft.

## Installation

    git clone github.com/Debaru/goutils/fileutil
    go get -u github.com/Debaru/goutils/fileutil
    go get -u github.com/mholt/archiver
    go get -u golang.org/x/net/html
    go build curscraft/cmd/curscraft/curscraft.go

## Usage

    -set-path [PATH] : Path to your World Of Warcraft Addons repository.
    -list : List of your World Of Warcraft Addons.
    -update-all : Update all your World Of Warcraft Addons.
    -help : Display help.

Configuration file are stored here :
* Linux : ~/.config/curscraft/config.txt
* Windows : HOMEDIR/AppData/Local/curscraft/config.txt
