package controllers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/viewmodels"
)

type configController struct {
	template *template.Template
}

func (c *configController) get(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: config :: get - %v\n", req.URL)

	w.Header().Set("Content-Type", "text/html")
	err := c.template.Execute(w, viewmodels.GetConfig())

	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Failed to execute the template", 500)
	}
}
