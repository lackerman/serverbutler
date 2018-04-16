package utils

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

	paths := new([]string)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && strings.IndexRune(fileInfo.Name(), '.') != 0 {
			*paths = append(*paths, fileInfo.Name())
		}
	}

	return *paths, nil
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
					read = read[0: len(read)-1]
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

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func DownloadFile(dir string, url string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	filename := filepath.Base(url)
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return err
	}
	defer file.Close()

	numCopied, err := io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	if numCopied == 0 {
		return errors.New("No data was downloaded from the URL")
	}
	return nil
}

// UnzipFile takes a destination folder unzips the supplied file to the destination folder
func UnzipFile(dst string, filename string) error {
	if !FileExists(dst) {
		if err := os.Mkdir(dst, os.ModePerm); err != nil {
			return err
		}
	}
	// Open a zip archive for reading.
	r, err := zip.OpenReader(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		nf, err := os.Create(filepath.Join(dst, f.Name))
		if err != nil {
			return err
		}
		_, err = io.Copy(nf, rc)
		if err != nil {
			return err
		}
		rc.Close()
		nf.Close()
	}
	return nil
}
