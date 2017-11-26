package controllers

import (
	"html/template"
	"net/http"

	"github.com/syndtr/goleveldb/leveldb"
)

// RegisterRoutes used by the web app to direct to controllers and static assest handlers
func RegisterRoutes(templates *template.Template, db *leveldb.DB) {
	home := &homeController{templates.Lookup("home.html")}
	config := &configController{templates.Lookup("config.html"), db}
	static := &staticController{"public/"}

	http.HandleFunc("/", home.get)
	http.HandleFunc("/config/", config.get)
	http.HandleFunc("/styles/", static.handler)
	http.HandleFunc("/scripts/", static.handler)

	ip := &ipController{}
	slack := &slackController{db}
	openvpn := &openvpnController{db}

	http.HandleFunc("/api/ip", ip.get)
	http.HandleFunc("/api/slack", slack.update)
	http.HandleFunc("/api/openvpn", openvpn.update)
}
