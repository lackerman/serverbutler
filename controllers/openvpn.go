package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/utils"
)

type openvpnController struct {
	template *template.Template
}

func (c *openvpnController) getConfigFiles(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: openvpnController :: getConfigFiles - %v\n", req.URL)

	paths, err := utils.GetFileList("controllers/")
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Failed to execute the template", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(paths)
}
