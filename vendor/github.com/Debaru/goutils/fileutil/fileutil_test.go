package fileutil

import (
	"fmt"
	"path/filepath"
	"testing"
)

// Testing GetName
func TestGetName(t *testing.T) {
	testCase := []struct {
		name string
		want string
	}{
		{"Test.txt", "Test"},
		{"test", "test"},
		{"Test A.txt", "Test A"},
		{"Test_A.txt", "Test_A"},
		{"Test.A.txt", "Test.A"},
		{"Test_A.B.txt", "Test_A.B"},
		{"Test_A.B..txt", "Test_A.B."},
		{"Test A..txt", "Test A."},
		{"Test A...txt", "Test A.."},
	}

	for _, tc := range testCase {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) {
			n := GetName(tc.name)
			if n != tc.want {
				t.Errorf("Expected %s, got %s", tc.want, n)
			}
		})
	}
}

// Testing GetExtension
func TestGetExtension(t *testing.T) {
	testCase := []struct {
		name string
		want string
	}{
		{"Test.txt", "txt"},
		{"test", ""},
		{"Test A.txt", "txt"},
		{"Test_A.txt", "txt"},
		{"Test.A.txt", "txt"},
		{"Test_A.B.txt", "txt"},
		{"Test_A.B..txt", "txt"},
		{"Test A..txt", "txt"},
		{"Test A...txt", "txt"},
	}

	for _, tc := range testCase {
		t.Run(fmt.Sprintf("%s", tc.name), func(t *testing.T) {
			e := GetExtension(tc.name)
			if e != tc.want {
				t.Errorf("Expected %s, got %s", tc.want, e)
			}
		})
	}
}

// Testing Download
func TestDownload(t *testing.T) {
	TestCase := []struct {
		url  string
		want string
		err  bool
	}{
		{"https://golang.org/doc/gopher/frontpage.png", "frontpage.png", false},
		{"https://golang.org/doc/gopher/", "", true},
		{"", "", true},
		{"https://golang.org/", "", true},
		{"https://golang.org", "golang.org", false},
	}

	for _, tc := range TestCase {
		t.Run(fmt.Sprintf("%s", tc.url), func(t *testing.T) {
			f, err := Download(tc.url, "test_data")

			if err != nil {
				if tc.err == false {
					t.Errorf("Error : Expected %t, got %s", tc.err, err)
				}
			}

			if err == nil {
				n := filepath.Base(f.Name())
				if n != tc.want {
					t.Errorf("Expected %s, got %s", tc.want, f.Name())
				}
			}
		})
	}
}

// Testing Copy
func TestCopy(t *testing.T) {
	TestCase := []struct {
		name string
		src  string
		dest string
		err  bool
	}{
		{"Normal", "test_data/test_copy.txt", "test_data/test_copy_2.txt", false},
		{"No_Src", "test_data/no_src.txt", "test_data/test_copy_2.txt", true},
		{"No_Dest", "test_data/test_copy.txt", "", true},
	}

	for _, tc := range TestCase {
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.src, tc.dest)

			if err != nil {
				if tc.err == false {
					t.Error("No error was expected")
				}
			}
		})
	}
}

// Testing CopyCut
func TestCopyCut(t *testing.T) {
	TestCase := []struct {
		name string
		src  string
		dest string
		err  bool
	}{
		{"Normal", "test_data/test_copy_2.txt", "test_data/test_cut.txt", false},
		{"No_Src", "test_data/no_src.txt", "test_data/test_cut.txt", true},
		{"No_Dest", "test_data/test_copy_2.txt", "", true},
	}

	for _, tc := range TestCase {
		t.Run(tc.name, func(t *testing.T) {
			err := CopyCut(tc.src, tc.dest)

			if err != nil {
				if tc.err == false {
					t.Error("No error was expected")
				}
			}
		})
	}
}
