// Package fileutil provide some utils to help with files
package fileutil

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetExtension return extension of file
func GetExtension(name string) string {
	i := strings.LastIndex(name, ".")
	if i < 0 {
		return ""
	}
	return name[i+1:]
}

// GetName return name of file without extension
func GetName(name string) string {
	i := strings.LastIndex(name, ".")
	if i < 0 {
		return name
	}
	return name[:i]
}

// Download url into directoy. Return os.File if method is successful
func Download(url string, dir string) (*os.File, error) {
	// File name
	_, fName := filepath.Split(url)
	fPath := filepath.Join(dir, fName)

	// Open URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create File
	f, err := os.Create(fPath)
	if err != nil {
		return nil, err
	}

	// Copy URL body to new file
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// Copy copies file from source (src) to destination (dst)
func Copy(src string, dest string) error {
	// Open source file
	fSrc, err := os.Open(src)
	defer fSrc.Close()
	if err != nil {
		return err
	}

	// Create destination file
	fDest, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0755)
	defer fDest.Close()
	if err != nil {
		return err
	}

	// Copy content
	_, err = io.Copy(fDest, fSrc)
	if err != nil {
		return err
	}

	return nil
}

// CopyCut copies file from source (src) to destination (dst), then
// delete source (src) file
func CopyCut(src string, dest string) error {
	// Copy
	err := Copy(src, dest)
	if err != nil {
		return err
	}

	// Delete source file
	err = os.Remove(src)
	if err != nil {
		return err
	}

	return nil
}
