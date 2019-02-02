package handlers

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"

	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/viewmodels"
)

func IpHandler(ctx *gin.Context) {
	client := &http.Client{
		Timeout:   3 * time.Second,
		Transport: &http.Transport{},
	}

	ipInfo, err := getIPInfo(client)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, ipInfo)
}

func getIPInfo(client *http.Client) (*viewmodels.IPInfo, error) {
	cmd := exec.Command("curl", "-vvv", "http://ipecho.net/plain")
	bites, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bites))

	res, err := client.Get("http://ipecho.net/plain")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("calling http://ipecho.net/plain was unsuccessful: %v", res.StatusCode)
	}

	// Get the IP
	body, err := ioutil.ReadAll(res.Body)
	ip := string(body)

	// Query the information using the IP
	res, err = client.Get(fmt.Sprintf("https://ipapi.co/%v/json", ip))
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("calling https://ipapi.co/%v/json was unsuccessful: %v", res.StatusCode)
	}

	// Return the IP Information from the previous client call
	ipInfo := &viewmodels.IPInfo{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(ipInfo); err != nil {
		return nil, err
	}
	return ipInfo, nil
}
