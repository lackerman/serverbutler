package main

import (
	"embed"
	"flag"
	"fmt"
	"net/http"

	"github.com/lackerman/serverbutler/handlers"
	"github.com/syndtr/goleveldb/leveldb"

	"k8s.io/klog/v2"
	"k8s.io/klog/v2/klogr"
)

//go:embed templates
var content embed.FS

func main() {
	port := flag.Int("port", 8080, "The port to use for the server (default: 8080)")
	path := flag.String("path", "./bin/tmp", "The file path to use for level db")
	flag.Parse()

	klog.InitFlags(nil)

	logger := klogr.New().WithName("serverbutler")
	logger.V(0).Info("starting app", "port", port, "path", path)

	// The returned DB instance is safe for concurrent use. Which means that all
	// DB's methods may be called concurrently from multiple goroutines.
	db, err := leveldb.OpenFile(*path, nil)
	if err != nil {
		logger.Error(err, "failed to open the connection to leveldb")
	}

	router := handlers.RegisterRoutes(logger, db, &content)
	if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), router); err != nil {
		logger.Error(err, "failed to start up the server")
	}
}
