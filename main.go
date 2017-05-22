package main

import (
	log "github.com/sirupsen/logrus"

	"runtime"

	"github.com/spf13/pflag"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var dataFolder string
	var indexPath string
	var addr string
	var forceRebuild bool
	var baseURL string
	pflag.StringVar(&dataFolder, "data-path", "", "Path to the pyvideo data folder")
	pflag.StringVar(&indexPath, "index-path", "search.bleve", "Path to the search index folder")
	pflag.StringVar(&addr, "http-addr", "127.0.0.1:8080", "Address the HTTP server should listen on for API calls")
	pflag.BoolVar(&forceRebuild, "force-rebuild", false, "Rebuild the index even if it already exists")
	pflag.StringVar(&baseURL, "base-url", "http://pyvideo.org", "Base URL of the pyvideo website")
	pflag.Parse()

	idx, err := loadIndex(indexPath, dataFolder, forceRebuild)
	if err != nil {
		log.WithError(err).Fatalf("Failed to load index on %s", indexPath)
	}
	defer idx.Close()

	if err := runHTTPD(idx, addr); err != nil {
		log.WithError(err).Fatalf("Failed to start HTTPD on %s", addr)
	}
}