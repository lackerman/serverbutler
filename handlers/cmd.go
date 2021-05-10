package handlers

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/lackerman/serverbutler/viewmodels"
)

type cmdHandler struct {
	template *template.Template
	logger   logr.Logger
}

func CmdHandler(t *template.Template, logger logr.Logger) *cmdHandler {
	return &cmdHandler{template: t, logger: logger}
}

func (c *cmdHandler) get(ctx *gin.Context) {
	c.template.Execute(ctx.Writer, viewmodels.Config{
		Site: viewmodels.Site{Page: "Command"},
	})
}

func (c *cmdHandler) execute(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	args := strings.Split(string(body), " ")
	cmd := exec.Command(args[0], args[1:]...)
	bites, err := cmd.CombinedOutput()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	if _, err = fmt.Fprintf(ctx.Writer, string(bites)); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}
