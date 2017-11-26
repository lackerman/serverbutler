package utils

import (
	"encoding/json"
	"log"
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
			log.Printf("Loading template: %v", path)
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
	payload, err := json.Marshal(ErrorMessage{msg})
	if err != nil {
		http.Error(w, "Failed to marshal error. Original message: "+msg, 500)
		return
	}
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(payload)
}
