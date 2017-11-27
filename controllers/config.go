package controllers

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/lackerman/serverbutler/constants"
	"github.com/lackerman/serverbutler/utils"
	"github.com/lackerman/serverbutler/viewmodels"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type configController struct {
	template *template.Template
	db       *leveldb.DB
}

func (c *configController) get(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: config :: get - %v\n", req.URL)

	slack, err := c.slackConfig()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	openvpn, err := c.openvpnConfig()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	config := viewmodels.Config{
		Title:   "ServerConf",
		Heading: "ServerConf",
		OpenVPN: *openvpn,
		Slack:   *slack,
	}

	w.Header().Set("Content-Type", "text/html")
	err = c.template.Execute(w, config)
	if err != nil {
		http.Error(w, "Failed to execute the template", 500)
	}
}

func (c *configController) openvpnConfig() (*viewmodels.OpenVPN, error) {
	log.Println("controller :: configController :: openvpnConfig")

	selected, err := c.db.Get([]byte(constants.OpenvpnSelected), nil)
	if err != nil {
		if err != errors.ErrNotFound { // Ignoring missing config at this point in time
			return nil, errors.New("Failed to retrieve openvpn config. Reason: " + err.Error())
		}
	}

	dirBytes, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		if err == errors.ErrNotFound {
			// Ignoring missing config at this point in time
			dirBytes = []byte("/opt/openvpn/")
			c.db.Put([]byte(constants.OpenvpnDir), dirBytes, nil)
		} else {
			return nil, errors.New("Failed to retrieve openvpn config. Reason: " + err.Error())
		}
	}
	dir := string(dirBytes)

	paths, err := utils.GetFileList(string(dir))
	if err != nil {
		return nil, errors.New("Failed to execute the template. Reason: " + err.Error())
	}
	sort.Strings(paths)

	contents, err := ioutil.ReadFile(dir + "/" + constants.OpenvpnCredentialFile)
	creds := []string{"", ""}
	if err == nil {
		creds = strings.Split(string(contents), "\n")
	}
	username := creds[0]
	password := utils.Hash(creds[1])

	return &viewmodels.OpenVPN{
		Configs:  paths,
		Selected: string(selected),
		Username: username,
		Password: password,
	}, nil
}

func (c *configController) slackConfig() (*viewmodels.Slack, error) {
	url, err := c.db.Get([]byte(constants.SlackURLKey), nil)
	if err != nil {
		return nil, errors.New("Failed to retrieve slack url. Reason: " + err.Error())
	}
	return &viewmodels.Slack{
		URL: string(url),
	}, nil
}
