package utils

import (
	"log"
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
