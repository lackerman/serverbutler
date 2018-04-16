package utils

import (
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

type BinaryFileSystem struct {
	fs http.FileSystem
}

func (b *BinaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *BinaryFileSystem) Exists(prefix string, filepath string) bool {

	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

// BinaryFileSystem is an in-memory representation of the go-bindata files
func NewBinaryFileSystem(root string) *BinaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root}
	return &BinaryFileSystem{fs}
}
