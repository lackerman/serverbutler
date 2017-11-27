package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
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
