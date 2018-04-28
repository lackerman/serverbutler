package handlers

import (
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"
)

// RegisterRoutes used by the web app to direct to handlers and static assest handlers
func RegisterRoutes(templates *template.Template, db *leveldb.DB) error {
	r := gin.Default()

	r.SetHTMLTemplate(templates)

	r.GET("/", HomeHandler("home.html"))
	r.GET("/ip", IpHandler)
	r.GET("/config", NewConfigHandler("config.html", db).get)
	r.POST("/api/slack/config", SlackHandler(db))

	oh := NewOpenvpnHandler(db)
	ovpn := r.Group("/api/openvpn")

	ovpn.POST("/config", oh.saveConfigDir)
	ovpn.POST("/selection", oh.saveSelection)
	ovpn.POST("/download", oh.downloadConfig)
	ovpn.POST("/restart", oh.restart)
	ovpn.GET("/credentials", oh.credentials)

	r.Run(":3000")
	return nil
}
