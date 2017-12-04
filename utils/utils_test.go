package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestSearchAndReplace(t *testing.T) {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	ignored := "file.dat"
	dir := "tmp"

	p, _ := filepath.Abs(dir)
	os.Mkdir(dir, os.ModePerm)
	// Setup tmp files for test
	for _, fn := range files {
		writeFile(t, p+"/"+fn, "hello")
	}
	writeFile(t, p+"/"+ignored, "hello")

	// Execute
	SearchAndReplace(p, "*.txt", "hello", "goodbye")

	// Validate
	for _, file := range files {
		contents := readFile(t, dir, file)
		if contents != "goodbye" {
			t.Errorf("The file '%v' was not updated", file)
		}
	}

	contents := readFile(t, dir, ignored)
	if contents != "hello" {
		t.Errorf("The Search and Replace was not contained to the file mask")
	}

	// Teardown
	err := os.RemoveAll(dir)
	if err != nil {
		t.Errorf("Failed to cleanup after the test")
	}
}

func readFile(t *testing.T, dir string, file string) string {
	contents, err := ioutil.ReadFile(dir + "/" + file)
	if err != nil {
		t.Errorf("Failed to read the file. %v", err)
	}
	return string(contents)
}

func writeFile(t *testing.T, path string, output string) {
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create the file: %v. %v", path, err.Error())
	}
	defer file.Close()
	_, err = file.WriteString(output)
	if err != nil {
		t.Fatalf("Failed to write contents to file. %v", err.Error())
	}
}
