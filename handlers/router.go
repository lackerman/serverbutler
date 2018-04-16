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
	r.GET("/config/", NewConfigHandler("config.html", db).get)
	r.POST("/api/slack/config", SlackHandler(db))

	openvpn := NewOpenvpnHandler(db)
	r.Group("/api/openvpn")

	r.POST("/config", openvpn.saveConfigDir)
	r.POST("/selection", openvpn.saveSelection)
	r.POST("/download", openvpn.downloadConfig)
	r.POST("/restart", openvpn.restart)
	r.GET("/credentials", openvpn.credentials)

	r.Run(":3000")
	return nil
}
