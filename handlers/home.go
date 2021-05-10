package handlers

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/viewmodels"
)

func HomeHandler(t *template.Template) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t.Execute(ctx.Writer, viewmodels.Site{Page: "Home"})
	}
}
