package handlers

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"fmt"

	"github.com/gin-gonic/gin"
)

func CmdHandler(ctx *gin.Context) {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
	args := strings.Split(string(body), " ")
	cmd := exec.Command(args[0], args[1:]...)
	bites, err := cmd.CombinedOutput()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}

	if _, err = fmt.Fprintf(ctx.Writer, string(bites)); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
	}
}
