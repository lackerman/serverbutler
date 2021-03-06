package handlers

import (
	"io/ioutil"
	"net/http"
	"time"

	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/viewmodels"
)

func IpHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		client := &http.Client{
			Timeout:   3 * time.Second,
			Transport: &http.Transport{},
		}

		ipInfo, err := getIPInfo(client)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, ipInfo)
	}
}

func getIPInfo(client *http.Client) (*viewmodels.IPInfo, error) {
	res, err := client.Get("http://ipecho.net/plain")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Get the response body
	bites, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("calling http://ipecho.net/plain was unsuccessful. failed to read response body: %+v", err)
	}
	ip := string(bites)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("calling http://ipecho.net/plain was unsuccessful. %d: %v", res.StatusCode, ip)
	}

	// Query the information using the IP
	res, err = client.Get(fmt.Sprintf("http://ip-api.com/json/%v", ip))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		b, err1 := ioutil.ReadAll(res.Body)
		if err1 != nil {
			return nil, fmt.Errorf("calling http://ip-api.com/json/%v was unsuccessful. %d: %+v", ip, res.StatusCode, err1)
		}
		return nil, fmt.Errorf("calling http://ip-api.com/json/%v was unsuccessful. %d: %+v", ip, res.StatusCode, string(b))
	}

	// Return the IP Information from the previous client call
	ipInfo := &viewmodels.IPInfo{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(ipInfo); err != nil {
		return nil, err
	}

	return ipInfo, nil
}
