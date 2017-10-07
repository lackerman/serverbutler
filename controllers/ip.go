package controllers

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
)

type ipController struct {
}

func (a *ipController) get(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: api :: get - %v\n", req.URL)

	ip, err := retrieveIP()
	if err != nil {
		http.Error(w, `{ "message": "Failed to retrieve current IP" }`, 500)
	}

	res, err := http.Get(fmt.Sprintf("http://ip-api.com/json/%v", ip))
	if err != nil {
		http.Error(w, `{ "message": "Failed to retrieve IP information" }`, 500)
	}

	w.Header().Add("Content-Type", "application/json")
	reader := bufio.NewReader(res.Body)
	reader.WriteTo(w)
}

func retrieveIP() (string, error) {
	res, err := http.Get("http://ipecho.net/plain")
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	return buf.String(), nil
}
