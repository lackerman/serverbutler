package handlers

import (
	"net/http"

	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/viewmodels"
)

func IpHandler(ctx *gin.Context) {
	ipInfo, err := getIPInfo()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ipInfo)
}

func getIPInfo() (*viewmodels.IPInfo, error) {
	res, err := http.Get("http://ipecho.net/plain")
	if err != nil {
		return nil, err
	}

	// Get the IP
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	ip := buf.String()

	res, err = http.Get(fmt.Sprintf("https://ipapi.co/%v/json", ip))
	if err != nil {
		return nil, err
	}

	// Return the IP Information from the previous client call
	ipInfo := &viewmodels.IPInfo{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(ipInfo); err != nil {
		return nil, err
	}
	return ipInfo, nil
}
