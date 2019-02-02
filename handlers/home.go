package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/viewmodels"
)

func HomeHandler(template string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		site := viewmodels.Site{Page: "Home"}
		ctx.HTML(http.StatusOK, template, viewmodels.Home{Site: site})
	}
}
