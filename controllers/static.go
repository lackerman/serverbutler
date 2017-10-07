package controllers

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type staticController struct {
	publicDirectory string
}

func (c *staticController) handler(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: static :: handler - %v\n", req.URL)

	path := c.publicDirectory + req.URL.Path
	file, err := os.Open(path)
	defer file.Close()
	if err == nil {
		w.Header().Set("Content-Type", contentType(path))
		bufferedReader := bufio.NewReader(file)
		bufferedReader.WriteTo(w)
	} else {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Failed to find the requested resource")
	}
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
