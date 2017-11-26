package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/utils"
)

type ipController struct {
}

func (a *ipController) get(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: api :: get - %v\n", req.URL)

	res, err := http.Get("http://ipecho.net/plain")
	if err != nil {
		utils.WriteJSONError(w, 500, "Failed to retrieve current IP")
		return
	}

	// Get the IP
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	ip := buf.String()

	res, err = http.Get(fmt.Sprintf("https://ipapi.co/%v/json", ip))
	if err != nil {
		utils.WriteJSONError(w, 500, "Failed to retrieve IP information")
		return
	}

	// Return the IP Information from the previous client call
	w.Header().Add("Content-Type", "application/json")
	reader := bufio.NewReader(res.Body)
	reader.WriteTo(w)
}
