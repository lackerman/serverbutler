package main

import (
	"flag"
	"fmt"
	"github.com/lackerman/serverbutler/handlers"
	"github.com/syndtr/goleveldb/leveldb"
	"html/template"
	"net/http"
)

func main() {
	port := flag.Int("port", 8080, "The port to use for the server (default: 8080)")
	path := flag.String("path", "./bin/tmp", "The file path to use for level db")
	flag.Parse()

	templates := template.Must(template.ParseGlob("templates/*"))
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
