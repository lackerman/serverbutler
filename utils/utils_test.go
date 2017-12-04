package utils

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestSearchAndReplace(t *testing.T) {
	dir := "./utils_test_files"
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	shouldNotBeEdited := "file.dat"

	SearchAndReplace(dir, "*.txt", "hello", "goodbye")

	for _, file := range files {
		contents := readRelativeFile(t, dir, file)
		if contents != "goodbye" {
			t.Errorf("The file '%v' was not updated", file)
		}
	}

	contents := readRelativeFile(t, dir, shouldNotBeEdited)
	if contents != "hello" {
		t.Errorf("The Search and Replace was not contained to the file mask")
	}

	// Put the files back the way they were
	SearchAndReplace(dir, "*.txt", "goodbye", "hello")
}

func readRelativeFile(t *testing.T, dir string, file string) string {
	p, _ := filepath.Abs(dir + "/" + file)
	contents, err := ioutil.ReadFile(p)
	if err != nil {
		t.Errorf("Failed to read the file. %v", err)
	}
	return string(contents)
}
