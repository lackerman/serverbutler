//go:generate go-bindata -pkg utils -o utils/assets.go public templates
package main

import (
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/constants"
	"github.com/lackerman/serverbutler/controllers"
	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
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
	initDb(db)

	controllers.RegisterRoutes(templates, db)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDb(db *leveldb.DB) {
	if _, err := db.Get([]byte(constants.OpenvpnDir), nil); err == errors.ErrNotFound {
		db.Put([]byte(constants.OpenvpnDir), []byte("tmp/config"), nil)
	}
	if _, err := db.Get([]byte(constants.SlackURLKey), nil); err == errors.ErrNotFound {
		db.Put([]byte(constants.SlackURLKey), []byte(""), nil)
	}
}
