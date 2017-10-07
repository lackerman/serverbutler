package main

import (
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/controllers"
	"github.com/lackerman/serverbutler/utils"
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

	controllers.RegisterRoutes(templates)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
