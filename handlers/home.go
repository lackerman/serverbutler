package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/viewmodels"
)

func HomeHandler(template string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ipInfo, err := getIPInfo()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		site := viewmodels.Site{
			Page:    "Home",
			Heading: "Serverbutler",
		}
		ctx.HTML(http.StatusOK, template, viewmodels.Home{Site: site, IPInfo: ipInfo})
	}
}
