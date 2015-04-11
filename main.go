package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/golang/glog"
)

const defaultConfigPath = "lanky.json"

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", defaultConfigPath, "path to lanky configuration file (default "+defaultConfigPath+")")
	flag.Parse()

	r, err := os.Open(configPath)
	if err != nil {
		glog.Fatalf("Unable to open config file %v: %v", configPath, err)
	}

	config := &Config{}
	err = LoadConfig(r, config)
	if err != nil {
		glog.Fatalf("Error reading config file %v: %v", configPath, err)
	}

	stats := NewStats()

	RegisterRoutes(config, stats)

	handler := &LoggingHandler{http.DefaultServeMux, stats}
	address := config.Address
	cert := config.CertificatePath
	key := config.KeyPath

	glog.Warningf("Starting server listening at %v.", config.Address)
	glog.Fatal(ListenAndServe(address, cert, key, handler))
}
