package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/lackerman/serverbutler/utils"
	"io/ioutil"
)

type staticController struct {
	publicDirectory string
}

func (c *staticController) handler(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: static :: handler - %v\n", req.URL)

	bytes, err := c.staticAsset(req.URL.Path)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintln(w, "Failed to find the requested resource")
	} else {
		w.Header().Set("Content-Type", contentType(req.URL.Path))
		w.Write(bytes)
	}
}

func (c *staticController) staticFile(p string) ([]byte, error) {
	path := c.publicDirectory + p
	return ioutil.ReadFile(path)
}

func (c *staticController) staticAsset(p string) ([]byte, error) {
	return utils.Asset(p[1:])
}

func contentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".js"):
		return "text/javascript"
	}
	return "text/plain"
}
