package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/lackerman/serverbutler/constants"
	"github.com/lackerman/serverbutler/viewmodels"
	"net/http"
)

func SlackHandler(db *leveldb.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.PostForm("webhook")
		err := db.Put([]byte(constants.SlackURLKey), []byte(url), nil)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, &viewmodels.ErrorMessage{Message: "Failed to save slack config"})
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, "/config")
	}
}
