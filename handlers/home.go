package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/viewmodels"
)

func HomeHandler(template string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ipInfo, err := getIpInfo()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.HTML(http.StatusOK, template, viewmodels.Home{
			Title:   "Home",
			Heading: "Serverbutler",
			IpInfo:  ipInfo,
		})
	}
}
