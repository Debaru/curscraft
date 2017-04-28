package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var path string

func init() {
	u, _ := user.Current()
	env := runtime.GOOS

	switch env {
	case "linux":
		// Creating curscraft repository if isn't exist
		path = filepath.Join(u.HomeDir, ".config", "curscraft")
		if _, err := os.Open(path); os.IsNotExist(err) {
			err = os.Mkdir(path, 0740)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		}
	}
}

// SaveFile saving file name into path
// in json format
func SaveFile(name string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	p := filepath.Join(path, name)
	err = ioutil.WriteFile(p, data, 0740)
	if err != nil {
		return err
	}

	return nil
}

// OpenFile opening p file
// and load content in v
func OpenFile(name string, v interface{}) error {
	p := filepath.Join(path, name)
	data, err := ioutil.ReadFile(p)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		err = json.Unmarshal(data, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// Fatal display error message
// And exit of programm
func Fatal(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
