//go:generate go-bindata -ignore .*.swp -pkg utils -o utils/assets.go public templates
package main

import (
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/controllers"
	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

func main() {
	templates, err := utils.ParseTemplatesFromBinData("templates")
	if err != nil {
		panic(err.Error())
	}

	// The returned DB instance is safe for concurrent use. Which means that all
	// DB's methods may be called concurrently from multiple goroutines.
	db, err := leveldb.OpenFile("db", nil)
	if err != nil {
		panic(err.Error())
	}

	controllers.RegisterRoutes(templates, db)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
