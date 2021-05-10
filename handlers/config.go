package handlers

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"

	"github.com/lackerman/serverbutler/constants"
	"github.com/lackerman/serverbutler/utils"
	"github.com/lackerman/serverbutler/viewmodels"
)

type configController struct {
	db       *leveldb.DB
	logger   logr.Logger
	template string
}

func NewConfigHandler(t string, db *leveldb.DB, logger logr.Logger) *configController {
	return &configController{
		db:       db,
		template: t,
		logger:   logger,
	}
}

func (c *configController) get(ctx *gin.Context) {
	slack, err := c.slack()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	openvpn, err := c.openvpn()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, c.template, viewmodels.Config{
		Site:    viewmodels.Site{Page: "Config"},
		OpenVPN: *openvpn,
		Slack:   *slack,
	})
}

func (c *configController) openvpn() (*viewmodels.OpenVPN, error) {
	c.logger.V(4).Info("openvpn - retrieving configs")

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

	c.logger.V(4).Info("Openvpn Dir: %+v", dir)
	if b != nil {
		configs, err = retrieveConfigs(dir)
		if err != nil {
			c.logger.V(4).Info("%+v", err)
			return nil, err
		}
		for _, cfg := range configs {
			c.logger.V(4).Info("%v\n", cfg)
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
	c.logger.V(4).Info("slack - getting config")
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
