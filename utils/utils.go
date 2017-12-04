package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetFileList returns all the files in the specified directory
func GetFileList(directory string) ([]string, error) {
	file, err := os.Open(directory)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var path string
	paths := new([]string)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			path = directory + fileInfo.Name()
			*paths = append(*paths, path)
		}
	}

	return *paths, nil
}

// ErrorMessage is the definition of a JSON error message
type ErrorMessage struct {
	Message string `json:"message"`
}

// WriteJSONError uses a ResponseWriter to write a json error message
func WriteJSONError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorMessage{msg})
}

// Hash to get a sha1 hash of a string
func Hash(s string) string {
	hasher := sha1.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

// SearchAndReplace scans through a directory and replaces a string in each of the files
func SearchAndReplace(path string, fileMask string, search string, replace string) error {
	return filepath.Walk(path, func(filePath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore directories
		if !fi.IsDir() {
			// File suffix should be externalised
			matched, err := filepath.Match(fileMask, fi.Name())
			if err != nil {
				return err
			}
			if matched {
				read, err := ioutil.ReadFile(filePath)
				if err != nil {
					return err
				}
				if read[len(read)-1] == 10 {
					read = read[0 : len(read)-1]
				}
				newContents := strings.Replace(string(read), search, replace, 1)
				fmt.Printf("|%s|%s|", read, newContents)
				err = ioutil.WriteFile(filePath, []byte(newContents), 0)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}
