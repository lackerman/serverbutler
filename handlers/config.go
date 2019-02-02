package handlers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"gitlab.com/lackerman/serverbutler/constants"
	"gitlab.com/lackerman/serverbutler/utils"
	"gitlab.com/lackerman/serverbutler/viewmodels"
)

type configController struct {
	db       *leveldb.DB
	logger   *log.Logger
	template string
}

func NewConfigHandler(t string, db *leveldb.DB) *configController {
	return &configController{
		db:       db,
		template: t,
		logger:   log.New(os.Stdout, "handler :: config - ", log.LstdFlags),
	}
}

func (c *configController) get(ctx *gin.Context) {
	slack, err := c.slack()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	openvpn, err := c.openvpn()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, c.template, viewmodels.Config{
		Site:    viewmodels.Site{Page: "Config"},
		OpenVPN: *openvpn,
		Slack:   *slack,
	})
}

func (c *configController) openvpn() (*viewmodels.OpenVPN, error) {
	c.logger.Println("openvpn - retrieving configs")

	var configs []string
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

func (c *configController) slack() (*viewmodels.Slack, error) {
	c.logger.Println("slack - getting config")
	url := ""
	b, err := c.db.Get([]byte(constants.SlackURLKey), nil)
	if err != nil {
		if err != errors.ErrNotFound {
			return nil, errors.New("Failed to retrieve slack url. Reason: " + err.Error())
		}
	} else {
		url = string(b)
	}
	return &viewmodels.Slack{URL: url}, nil
}
