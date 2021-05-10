package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/lackerman/serverbutler/constants"
	"github.com/lackerman/serverbutler/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

type openvpnHandler struct {
	db     *leveldb.DB
	logger logr.Logger
}

func NewOpenvpnHandler(db *leveldb.DB, logger logr.Logger) *openvpnHandler {
	return &openvpnHandler{db: db, logger: logger}
}

func (c *openvpnHandler) saveConfigDir(ctx *gin.Context) {
	dir := ctx.PostForm("dir")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = c.db.Put([]byte(constants.OpenvpnDir), []byte(dir), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/config")
}

func (c *openvpnHandler) saveSelection(ctx *gin.Context) {
	selected := ctx.PostForm("config")
	err := c.db.Put([]byte(constants.OpenvpnSelected), []byte(selected), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusFound, "/config")
}

func (c *openvpnHandler) downloadConfig(ctx *gin.Context) {
	b, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	dir := string(b)
	err = utils.DownloadFile(dir, "https://downloads.nordcdn.com/configs/archives/servers/ovpn.zip")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = utils.UnzipFile(dir, filepath.Join(dir, "zip"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = os.Remove(filepath.Join(dir, "zip"))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (c *openvpnHandler) credentials(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	data := []byte(fmt.Sprintf("%v\n%v", username, password))

	dir, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = ioutil.WriteFile(string(dir)+"/"+constants.OpenvpnCredentialFile, data, 0644)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, "/config")
}

func (c *openvpnHandler) restart(ctx *gin.Context) {
	err := func(ctx *gin.Context) error {
		configDir, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
		if err != nil {
			return err
		}
		selection, err := c.db.Get([]byte(constants.OpenvpnSelected), nil)
		if err != nil {
			return err
		}
		config := string(configDir) + "/" + string(selection)
		c.logger.V(4).Info("Starting openvpn with the selected config", config)
		cmd := exec.Command("openvpn", config)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return err
		}
		c.logger.V(4).Info("Started subprocess: ", cmd.Process.Pid)
		return nil
	}(ctx)

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/config")
}
