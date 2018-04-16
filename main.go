//go:generate go-bindata -ignore ".*.swp|.DS_Store" -pkg utils -o utils/assets.go templates
package main

import (
	"os"

	"github.com/lackerman/serverbutler/handlers"
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
	db, err := leveldb.OpenFile(os.Getenv("config_dir"), nil)
	if err != nil {
		panic(err.Error())
	}

	if err := handlers.RegisterRoutes(templates, db); err != nil {
		panic(err.Error())
	}
}
