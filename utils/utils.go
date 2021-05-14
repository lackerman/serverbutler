package utils

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetFileList returns all the files in the specified directory
func GetFileList(directory string) ([]string, error) {
	paths := []string{}
	err := filepath.Walk(directory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				paths = append(paths, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}
	return paths, nil
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

func DownloadFile(dir string, url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	filename := filepath.Base(url)
	file, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	numCopied, err := io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}
	if numCopied == 0 {
		return "", errors.New("no data was downloaded from the URL")
	}
	return filename, nil
}

// UnzipFile takes a destination folder unzips the supplied file to the destination folder
func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, f.Name)

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), 0777)
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
