package controllers

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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

	configs := []string{}
	selected, username, password, dir := "", "", "", ""

	b, err := c.db.Get([]byte(constants.OpenvpnSelected), nil)
	if err != nil {
		// Ignoring missing config at this point in time
		if err != errors.ErrNotFound {
			return nil, errors.New("Failed to retrieve openvpn config. Reason: " + err.Error())
		}
	} else {
		selected = string(b)
	}

	b, err = c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		// Ignoring missing config at this point in time (during 1st initialisation)
		if err != errors.ErrNotFound {
			return nil, errors.New("Failed to retrieve openvpn config. Reason: " + err.Error())
		}
	} else {
		dir = string(b)
	}

	if b != nil {
		configs, err = retrieveConfigs(dir)
		if err != nil {
			return nil, err
		}
		username, password = retrieveCredentials(dir)
	}

	return &viewmodels.OpenVPN{
		ConfigDir: dir,
		Configs:   configs,
		Selected:  selected,
		Username:  username,
		Password:  password,
	}, nil
}

func retrieveConfigs(dir string) ([]string, error) {
	paths, err := utils.GetFileList(dir)
	if err != nil {
		return nil, errors.New("Failed to execute the template. Reason: " + err.Error())
	}
	sort.Strings(paths)
	return paths, nil
}

func retrieveCredentials(dir string) (string, string) {
	contents, err := ioutil.ReadFile(filepath.Join(dir, constants.OpenvpnCredentialFile))
	if err != nil {
		return "", ""
	}
	creds := strings.Split(string(contents), "\n")
	return creds[0], utils.Hash(creds[1])
}

func (c *configController) slackConfig() (*viewmodels.Slack, error) {
	url := ""
	b, err := c.db.Get([]byte(constants.SlackURLKey), nil)
	if err != nil {
		if err != errors.ErrNotFound {
			return nil, errors.New("Failed to retrieve slack url. Reason: " + err.Error())
		}
	} else {
		url = string(b)
	}
	return &viewmodels.Slack{url}, nil
}
