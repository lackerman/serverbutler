package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type openvpnController struct {
	db *leveldb.DB
}

func (c *openvpnController) handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		c.update(w, req)
		return
	default:
		http.Error(w, "Unsupported request type", http.StatusBadRequest)
	}
}

func (c *openvpnController) update(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: openvpnController :: getConfigFiles - %v\n", req.URL)

	paths, err := utils.GetFileList("controllers/")
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Failed to execute the template", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paths)
}
