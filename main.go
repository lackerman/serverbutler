//go:generate go-bindata -ignore ".*.swp|.DS_Store" -pkg utils -o utils/assets.go templates
package main

import (
	"flag"
	"fmt"
	"github.com/lackerman/serverbutler/handlers"
	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "The port to use for the server (default: 8080)")
	path := flag.String("path", "./tmp", "The file path to use for level db")
	flag.Parse()

	templates, err := utils.ParseTemplatesFromBinData("templates")
	if err != nil {
		panic(err.Error())
	}
	// The returned DB instance is safe for concurrent use. Which means that all
	// DB's methods may be called concurrently from multiple goroutines.
	db, err := leveldb.OpenFile(*path, nil)
	if err != nil {
		panic(err.Error())
	}

	router := handlers.RegisterRoutes(templates, db)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		panic(err.Error())
	}
}
