package controllers

import (
	"log"
	"net/http"

	"github.com/lackerman/serverbutler/utils"

	"github.com/lackerman/serverbutler/constants"

	"github.com/syndtr/goleveldb/leveldb"
)

type slackController struct {
	db *leveldb.DB
}

func (c *slackController) handle(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		c.update(w, req)
		return
	default:
		http.Error(w, "Unsupported request type", http.StatusBadRequest)
	}
}

func (c *slackController) update(w http.ResponseWriter, req *http.Request) {
	log.Printf("controller :: slack :: update - %v\n", req.URL)

	req.ParseForm()
	url := req.Form.Get("webhook")
	err := c.db.Put([]byte(constants.SlackURLKey), []byte(url), nil)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to save slack config")
		return
	}

	http.Redirect(w, req, "/config", 301)
}
