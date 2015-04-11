package main

import (
	"net/http"
	_ "net/http/pprof"
	"strings"

	"github.com/golang/glog"
)

type ByteWriter struct {
	http.ResponseWriter
	Wrote  int
	Status int
}

func (bw *ByteWriter) Write(b []byte) (int, error) {
	i, err := bw.ResponseWriter.Write(b)
	bw.Wrote += i
	return i, err
}

func (bw *ByteWriter) WriteHeader(code int) {
	bw.Status = code
	bw.ResponseWriter.WriteHeader(code)
}

type LoggingHandler struct {
	http.Handler
	*RuntimeStats
}

func (lh *LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// AFAICT WriteHeader is only called when not 200 OK.
	bw := &ByteWriter{w, 0, 200}

	lh.Handler.ServeHTTP(bw, r)

	ip := strings.Split(r.RemoteAddr, ":")[0]
	if r.URL.Path != "/status" {
		err := lh.RuntimeStats.IncStatus(bw.Status)
		if err != nil {
			glog.Error(err.Error())
		}
	}
	glog.Infof("%v %v - \"%v %v %v\" %v %v", ip, r.Header.Get("User-Agent"), r.Method, r.URL.Path, r.Proto, bw.Status, bw.Wrote)
}

func HandleFuncConfig(path string, fn func(w http.ResponseWriter, r *http.Request, c *Config) error, c *Config) {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r, c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func RegisterRoutes(config *Config, stats *RuntimeStats) {
	// GitHub Post-Receive requests
	HandleFuncConfig("/_github", githubHandler, config)
	// Hubot API
	HandleFuncConfig("/_hubot", hubotHandler, config)
	// Jenkins callback
	HandleFuncConfig("/_builder", builderHandler, config)

	// Organisations repository listing
	HandleFuncConfig("/repositories", repositoryHandler, config)

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, config, stats)
	})

	// landing page
	HandleFuncConfig("/", rootHandler, config)
}

func ListenAndServe(address, cert, key string, handler http.Handler) (err error) {
	if key != "" && cert != "" {
		err = http.ListenAndServeTLS(address, cert, key, handler)
	} else {
		err = http.ListenAndServe(address, handler)
	}

	return
}
