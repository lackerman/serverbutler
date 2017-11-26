package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type slackController struct {
	db *leveldb.DB
}

func (c *slackController) handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		c.update(w, req)
		return
	default:
		http.Error(w, "Unsupported request type", http.StatusBadRequest)
	}
}

func (c *slackController) update(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: slackController :: getConfigFiles - %v\n", req.URL)

	paths, err := utils.GetFileList("controllers/")
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Failed to execute the template", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paths)
}
