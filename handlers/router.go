package handlers

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

// RegisterRoutes used by the web app to direct to handlers and static assest handlers
func RegisterRoutes(prefix string, templates *template.Template, db *leveldb.DB) error {
	r := gin.Default()

	r.SetHTMLTemplate(templates)

	p := r.Group("/" + prefix)

	p.GET("/", HomeHandler("home.html"))
	p.GET("/ip", IpHandler)
	p.GET("/config", NewConfigHandler("config.html", db).get)
	p.POST("/api/slack/config", SlackHandler(db))

	oh := NewOpenvpnHandler(db)
	ovpn := p.Group("/api/openvpn")

	ovpn.POST("/config", oh.saveConfigDir)
	ovpn.POST("/selection", oh.saveSelection)
	ovpn.POST("/download", oh.downloadConfig)
	ovpn.POST("/restart", oh.restart)
	ovpn.GET("/credentials", oh.credentials)

	r.Run(":3000")
	return nil
}
