package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	tmpDir = "tmp"
)

func TestSearchAndReplace(t *testing.T) {
	// Setup
	files, ignored := setup(t)

	// Execute
	SearchAndReplace(tmpDir, "*.txt", "hello", "goodbye")

	// Validate
	for _, file := range files {
		contents := readFile(t, file)
		if contents != "goodbye" {
			t.Errorf("The file '%v' was not updated", file)
		}
	}

	contents := readFile(t, ignored)
	if contents != "hello" {
		t.Errorf("The Search and Replace was not contained to the file mask")
	}

	// Teardown
	err := os.RemoveAll(tmpDir)
	if err != nil {
		t.Errorf("Failed to cleanup after the test")
	}
}

func setup(t *testing.T) ([]string, string) {
	p, _ := filepath.Abs(tmpDir)
	if err := os.Mkdir(tmpDir, os.ModePerm); err != nil {
		t.Fatalf("Failed to create temp directory. %v", err.Error())
	}
	files := []string{}
	for i := 1; i <= 3; i++ {
		file := filepath.Join(p, fmt.Sprintf("file%i.txt", i))
		files = append(files, file)
		writeFile(t, file, "hello")
	}
	ignored := filepath.Join(p, "file.dat")
	writeFile(t, ignored, "hello")
	return files, ignored
}

func readFile(t *testing.T, file string) string {
	contents, err := ioutil.ReadFile(file)
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
