package handlers

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/gin-gonic/gin"
	"log"
	"github.com/lackerman/serverbutler/utils"
	"github.com/lackerman/serverbutler/constants"
	"os"
	"net/http"
	"github.com/lackerman/serverbutler/viewmodels"
	"path/filepath"
	"fmt"
	"os/exec"
	"io/ioutil"
)

type openvpnHandler struct {
	db     *leveldb.DB
	logger *log.Logger
}

func NewOpenvpnHandler(db *leveldb.DB) *openvpnHandler {
	return &openvpnHandler{db: db, logger: log.New(os.Stdout, "handler :: openvpn - ", log.LstdFlags)}
}

func (c *openvpnHandler) saveConfigDir(ctx *gin.Context) {
	dir := ctx.PostForm("dir")
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		c.reportError(ctx, err, "Failed to create dir.")
		return
	}

	err = c.db.Put([]byte(constants.OpenvpnDir), []byte(dir), nil)
	if err != nil {
		c.reportError(ctx, err, "Failed to save config dir")
		return
	}
	ctx.Redirect(http.StatusFound, "/config")
}

func (c *openvpnHandler) saveSelection(ctx *gin.Context) {
	selected := ctx.PostForm("config")
	err := c.db.Put([]byte(constants.OpenvpnSelected), []byte(selected), nil)
	if err != nil {
		c.reportError(ctx, err, "Failed to save saveSelection")
		return
	}
	ctx.Redirect(http.StatusFound, "/config")
}

func (c *openvpnHandler) downloadConfig(ctx *gin.Context) {
	b, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		c.reportError(ctx, err, "Failed to retrieve config")
		return
	}
	dir := string(b)
	err = utils.DownloadFile(dir, "https://nordvpn.com/api/files/zip")
	if err != nil {
		c.reportError(ctx, err, "Failed to download config")
		return
	}
	err = utils.UnzipFile(dir, filepath.Join(dir, "zip"))
	if err != nil {
		c.reportError(ctx, err, "Failed to unzip the downloaded config")
		return
	}
	err = os.Remove(filepath.Join(dir, "zip"))
	if err != nil {
		c.reportError(ctx, err, "Failed to cleanup after unzipping")
		return
	}
}

func (c *openvpnHandler) credentials(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	data := []byte(fmt.Sprintf("%v\n%v", username, password))

	dir, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		c.reportError(ctx, err, "Failed to retrieve config")
		return
	}
	err = ioutil.WriteFile(string(dir)+"/"+constants.OpenvpnCredentialFile, data, 0644)
	if err != nil {
		c.reportError(ctx, err, "Failed to write credentials")
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
		c.logger.Printf("Starting openvpn with the selected config: %v", config)
		cmd := exec.Command("openvpn", config)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			return err
		}
		c.logger.Printf("Started subprocess %d.", cmd.Process.Pid)
		return nil
	}(ctx)

	if err != nil {
		c.reportError(ctx, err, err.Error())
		return
	}
	ctx.Redirect(http.StatusTemporaryRedirect, "/config")
}

func (c *openvpnHandler) reportError(ctx *gin.Context, err error, m string) {
	msg := fmt.Sprintf("%v. Error: %v", m, err.Error())
	c.logger.Printf(msg)
	ctx.JSON(http.StatusInternalServerError, &viewmodels.ErrorMessage{m})
}
