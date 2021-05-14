package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

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
	payload, err := getJsonPayload(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	dir := payload["dir"]
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		c.logger.Error(err, err.Error()+" -- "+dir)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = c.db.Put([]byte(constants.OpenvpnDir), []byte(dir), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Saved the config directory"})
}

func (c *openvpnHandler) saveSelection(ctx *gin.Context) {
	request, err := getJsonPayload(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	selected := request["selected"]
	err = makeSymlink(selected)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = c.db.Put([]byte(constants.OpenvpnSelected), []byte(selected), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = c.restartOpenvpn(selected)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Saved selection: %v", selected)})
}

func (c *openvpnHandler) downloadConfig(ctx *gin.Context) {
	dirBytes, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	vpnDir := string(dirBytes)
	filename, err := utils.DownloadFile(vpnDir, "https://downloads.nordcdn.com/configs/archives/servers/ovpn.zip")
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = utils.Unzip(filepath.Join(vpnDir, filename), vpnDir)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	err = os.Remove(filepath.Join(vpnDir, filename))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully downloaded the config"})
}

func (c *openvpnHandler) credentials(ctx *gin.Context) {
	request, err := getJsonPayload(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	data := []byte(fmt.Sprintf("%v\n%v", request["username"], request["password"]))

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
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully updated the credentials"})
}

func (c *openvpnHandler) restart(ctx *gin.Context) {
	configDir, err := c.db.Get([]byte(constants.OpenvpnDir), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	selection, err := c.db.Get([]byte(constants.OpenvpnSelected), nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	config := string(configDir) + "/" + string(selection)
	err = c.restartOpenvpn(config)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully restarted openvpn"})
}
func (c *openvpnHandler) restartOpenvpn(config string) error {
	c.logger.V(4).Info("restarting openvpn with the selected config", config)
	cmd := exec.Command("service", "openvpn", "restart")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	c.logger.V(4).Info("Started subprocess: ", cmd.Process.Pid)
	return nil
}

func makeSymlink(selectedConfig string) error {
	openvpnConf, err := openvpnCfg()
	if err != nil {
		return err
	}
	if _, err := os.Lstat(openvpnConf); err == nil {
		err = os.Remove(openvpnConf)
		if err != nil {
			return err
		}
	}
	return os.Symlink(selectedConfig, openvpnConf)
}

func openvpnCfg() (string, error) {
	path := "/usr/local/etc/openvpn/openvpn.conf"
	if runtime.GOOS != "freebsd" {
		dir, err := currPath()
		if err != nil {
			return "", err
		}
		path = dir + "/openvpn.conf"
	}
	return path, nil
}

func currPath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", nil
	}
	return dir, nil
}
