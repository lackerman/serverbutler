package handlers

import (
	"embed"
	"fmt"
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
				internalError := fmt.Sprintf("Interval Server Error. See error logs for timestamp %v", stamp)
				c.JSON(http.StatusInternalServerError, gin.H{"errors": []string{internalError}})
				tmpLogger.V(0).Info("failed request", "error_timestamp", stamp)
				return
			}
			if len(c.Errors) > 0 {
				tmpLogger.V(0).Info("known failure", "errors", c.Errors.Errors())
				c.JSON(c.Writer.Status(), gin.H{"errors": c.Errors.Errors()})
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
	root.Group("/").
		GET("/", home).
		GET("/ip", ip).
		GET("/command", cmd.get).
		GET("/config", cfg.get)

	api := root.Group("/api")
	api.Group("/openvpn").
		POST("/config", vpn.saveConfigDir).
		POST("/selection", vpn.saveSelection).
		POST("/download", vpn.downloadConfig).
		POST("/restart", vpn.restart).
		POST("/credentials", vpn.credentials)

	api.Group("/transmission").
		POST("/state", vpn.saveConfigDir).
		POST("/url", vpn.saveSelection)

	api.Group("/slack").
		POST("/config", slack)

	api.Group("/command").
		POST("execute", cmd.execute)

	return r
}
