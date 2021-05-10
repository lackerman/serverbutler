package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lackerman/serverbutler/constants"
	"github.com/syndtr/goleveldb/leveldb"
)

func SlackHandler(db *leveldb.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.PostForm("webhook")
		err := db.Put([]byte(constants.SlackURLKey), []byte(url), nil)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, "/config")
	}
}
