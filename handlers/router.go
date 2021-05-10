package handlers

import (
	"net/http"
	"time"

	"github.com/lackerman/serverbutler/constants"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/go-logr/logr"
)

// RegisterRoutes used by the web app to configure handlers to paths
func RegisterRoutes(logger logr.Logger, db *leveldb.DB) http.Handler {
	logger = logger.WithName("router")

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
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

	root := r.Group("/" + constants.SitePrefix())

	p := root.Group("/")
	p.GET("/", HomeHandler("home.html"))
	p.GET("/ip", IpHandler)
	p.POST("/cmd", CmdHandler)
	p.GET("/config", NewConfigHandler("config.html", db, logger).get)
	p.POST("/api/slack/config", SlackHandler(db))

	oh := NewOpenvpnHandler(db, logger)
	o := root.Group("/api/openvpn")
	o.POST("/config", oh.saveConfigDir)
	o.POST("/selection", oh.saveSelection)
	o.POST("/download", oh.downloadConfig)
	o.POST("/restart", oh.restart)
	o.GET("/credentials", oh.credentials)

	return r
}
