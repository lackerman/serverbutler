package handlers

import (
	"net/http"

	"github.com/lackerman/serverbutler/viewmodels"
	"github.com/gin-gonic/gin"
	"bytes"
	"fmt"
	"encoding/json"
)

func IpHandler(ctx *gin.Context) {
	ipInfo, err := getIpInfo()
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ipInfo)
}

func getIpInfo() (*viewmodels.IpInfo, error) {
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
	ipInfo := &viewmodels.IpInfo{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(ipInfo); err != nil {
		return nil, err
	}
	return ipInfo, nil
}
