package curscraft

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	DirName    string
	FileName   string
	Properties map[string]interface{}
}

func checkDirConfigIsExist(d string) error {
	if _, err := os.Lstat(d); err != nil {
		err = os.Mkdir(d, 0760)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) setDirConfig() (string, error) {
	var p string

	u, _ := user.Current()
	env := runtime.GOOS

	switch env {
	case "windows":
		p = filepath.Join(u.HomeDir, "AppData", "Local", c.DirName)
		err := checkDirConfigIsExist(p)
		if err != nil {
			return p, err
		}
		return p, nil

	case "linux":
		p = filepath.Join(u.HomeDir, ".config", c.DirName)
		err := checkDirConfigIsExist(p)
		if err != nil {
			return p, err
		}
		return p, nil
	}

	return p, errors.New("OS missing")
}

func (c *Config) Save() error {
	var t string

	p, err := c.setDirConfig()
	if err != nil {
		return err
	}

	p = filepath.Join(p, c.FileName)
	f, err := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	for i, v := range c.Properties {
		t = fmt.Sprintf("%v = %v \n", i, v)
		f.WriteString(t)
	}

	return nil
}

func (c *Config) Read() {
	var p map[string]interface{}
	p = make(map[string]interface{})

	r, _ := c.setDirConfig()

	r = filepath.Join(r, c.FileName)
	f, _ := os.Open(r)
	defer f.Close()

	s := bufio.NewScanner(f)
	for s.Scan() {
		t := s.Text()
		z := strings.SplitN(t, "=", 2)
		p[strings.TrimSpace(z[0])] = strings.TrimSpace(z[1])
	}

	c.Properties = p
}
