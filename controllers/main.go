package controllers

import (
	"html/template"
	"net/http"
)

// RegisterRoutes used by the web app to direct to controllers and static assest handlers
func RegisterRoutes(templates *template.Template) {
	home := new(homeController)
	home.template = templates.Lookup("home.html")

	config := new(configController)
	config.template = templates.Lookup("config.html")

	ip := new(ipController)
	openvpn := new(openvpnController)

	static := new(staticController)
	static.publicDirectory = "public/"

	http.HandleFunc("/", home.get)
	http.HandleFunc("/config/", config.get)
	http.HandleFunc("/styles/", static.handler)
	http.HandleFunc("/scripts/", static.handler)

	http.HandleFunc("/api/ip", ip.get)
	http.HandleFunc("/api/openvpn", openvpn.getConfigFiles)
}
