package handlers

import (
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/lackerman/serverbutler/constants"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/go-logr/logr"
)

// RegisterRoutes used by the web app to configure handlers to paths
func RegisterRoutes(logger logr.Logger, db *leveldb.DB, content *embed.FS) http.Handler {
	logger = logger.WithName("router")

	r := gin.Default()
	// Custom recovery and logging middleware
	r.Use(func(c *gin.Context) {
		t := time.Now()
		defer func(start time.Time) {
			latency := time.Since(start)
			tmpLogger := logger.WithValues(
				"method", c.Request.Method,
				"url", c.Request.URL.String(),
				"status", c.Writer.Status(),
				"latency_ms", latency.Milliseconds())

			if err := recover(); err != nil {
				stamp := time.Now().UnixNano()
				c.String(500, "Interval Server Error. See error logs for timestamp %v", stamp)
				tmpLogger.V(0).Info("failed request", "error_timestamp", stamp)
				return
			}
			if len(c.Errors) > 0 {
				tmpLogger.V(0).Info("known failure", "errors", c.Errors.Errors())
				return
			}
			tmpLogger.V(5).Info("successful request")
		}(t)
		c.Next()
	})

	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))

	home := HomeHandler(tmpl.Lookup("home.html"))
	cmd := CmdHandler(tmpl.Lookup("command.html"), logger)
	cfg := NewConfigHandler(tmpl.Lookup("config.html"), db, logger)
	vpn := NewOpenvpnHandler(db, logger)
	slack := SlackHandler(db)
	ip := IpHandler()

	root := r.Group("/" + constants.SitePrefix())
	p := root.Group("/")
	p.GET("/", home)
	p.GET("/ip", ip)
	p.GET("/command", cmd.get)
	p.GET("/config", cfg.get)
	p.POST("/api/slack/config", slack)
	p.POST("/api/command/execute", cmd.execute)

	o := root.Group("/api/openvpn")
	o.POST("/config", vpn.saveConfigDir)
	o.POST("/selection", vpn.saveSelection)
	o.POST("/download", vpn.downloadConfig)
	o.POST("/restart", vpn.restart)
	o.GET("/credentials", vpn.credentials)

	return r
}
