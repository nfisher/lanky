package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

const defaultConfigPath = "lanky.json"

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "lanky.json", "path to lanky configuration file (default "+defaultConfigPath+")")
	flag.Parse()

	r, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Unable to open config file %v: %v", configPath, err)
	}

	config := &Config{}
	err = LoadConfig(r, config)
	if err != nil {
		log.Fatalf("Error reading config file %v: %v", configPath, err)
	}

	serve(config)
}

func serve(config *Config) {
	// Janky paths
	// /_github - Post-Receive requests
	// /_builder - Jenkins callback
	// /_hubot - Hubot api
	// / - landing page

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootHandler(w, r, config)
	})

	log.Printf("Starting server listening at %v.", config.Address)
	log.Fatal(http.ListenAndServe(config.Address, nil))
}
