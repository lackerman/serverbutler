package handlers

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

// RegisterRoutes used by the web app to configure handlers to paths
func RegisterRoutes(templates *template.Template, db *leveldb.DB) http.Handler {
	r := gin.Default()

	r.SetHTMLTemplate(templates)

	p := r.Group("/")
	p.GET("/", HomeHandler("home.html"))
	p.GET("/ip", IpHandler)
	p.POST("/cmd", CmdHandler)
	p.GET("/config", NewConfigHandler("config.html", db).get)
	p.POST("/api/slack/config", SlackHandler(db))

	oh := NewOpenvpnHandler(db)
	o := p.Group("/api/openvpn")
	o.POST("/config", oh.saveConfigDir)
	o.POST("/selection", oh.saveSelection)
	o.POST("/download", oh.downloadConfig)
	o.POST("/restart", oh.restart)
	o.GET("/credentials", oh.credentials)

	return r
}
