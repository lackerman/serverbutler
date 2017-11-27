package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/constants"

	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type openvpnController struct {
	db *leveldb.DB
}

func (c *openvpnController) selection(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: openvpn :: selection - %v\n", req.URL)
	if req.Method != "POST" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Unsupported request type")
		return
	}
	req.ParseForm()
	selected := req.Form.Get("config")
	err := c.db.Put([]byte(constants.OpenvpnSelected), []byte(selected), nil)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to save selection")
		return
	}
	http.Redirect(w, req, "/config", 301)
}

func (c *openvpnController) credentials(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: openvpn :: credentials - %v\n", req.URL)
	if req.Method != "POST" {
		utils.WriteJSONError(w, http.StatusBadRequest, "Unsupported request type")
		return
	}

	req.ParseForm()
	username := req.Form.Get("username")
	password := req.Form.Get("password")
	data := []byte(fmt.Sprintf("%v\n%v", username, password))

	dir, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to retrieve config")
		return
	}
	err = ioutil.WriteFile(string(dir)+"/"+constants.OpenvpnCredentialFile, data, 0644)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to write credentials")
		return
	}
	http.Redirect(w, req, "/config", 301)
}
