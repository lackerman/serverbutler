package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lackerman/serverbutler/constants"

	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type openvpnController struct {
	db *leveldb.DB
}

func (c *openvpnController) updateConfigDir(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: openvpn :: config - %v\n", req.URL)
	if req.Method != "POST" {
		log.Printf("Incorrect request method: %v", req.Method)
		utils.WriteJSONError(w, http.StatusBadRequest, "Unsupported request type")
		return
	}
	req.ParseForm()
	dir := req.Form.Get("dir")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Printf(err.Error())
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to create dir.")
	}

	err = c.db.Put([]byte(constants.OpenvpnDir), []byte(dir), nil)
	if err != nil {
		log.Printf(err.Error())
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to save config dir")
	}
	http.Redirect(w, req, "/config", 301)
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

func (c *openvpnController) downloadConfig(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: openvpn :: credentials - %v\n", req.URL)
	b, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to retrieve config")
		return
	}
	dir := string(b)
	err = utils.DownloadFile(dir, "https://nordvpn.com/api/files/zip")
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to download config")
	}
	err = utils.UnzipFile(dir, filepath.Join(dir, "zip"))
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to unzip the downloaded config")
	}
	err = os.Remove(filepath.Join(dir, "zip"))
	if err != nil {
		log.Printf("Failed to delete the zip file. %v", err.Error())
	}
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

func (c *openvpnController) restart(w http.ResponseWriter, req *http.Request) {
	err := func (w http.ResponseWriter, req *http.Request) error {
		configDir, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
		if err != nil {
			return err
		}
		selection, err := c.db.Get([]byte(constants.OpenvpnSelected), nil)
		if err != nil {
			return err
		}
		config := string(configDir) + "/" + string(selection)
		log.Printf("Starting openvpn with the selected config: %v", config)
		cmd := exec.Command("openvpn", config)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return err
		}
		log.Printf("Started subprocess %d.", cmd.Process.Pid)
		return nil
	}(w, req)

	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, req, "/config", 301)
}