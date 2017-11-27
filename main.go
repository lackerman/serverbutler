package main

import (
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/controllers"
	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	templatePaths, err := utils.GetFileList("templates/")
	if err != nil {
		panic(err.Error())
	}
	templates, err := utils.ParseTemplates(templatePaths)
	if err != nil {
		panic(err.Error())
	}

	// The returned DB instance is safe for concurrent use. Which mean that all
	// DB's methods may be called concurrently from multiple goroutine.
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic(err.Error())
	}

	controllers.RegisterRoutes(templates, db)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
