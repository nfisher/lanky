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

func Serve(config *Config) {
	// GitHub Post-Receive requests
	http.HandleFunc("/_github", func(w http.ResponseWriter, r *http.Request) {
		githubHandler(w, r, config)
	})

	// Hubot API
	http.HandleFunc("/_hubot", func(w http.ResponseWriter, r *http.Request) {
		hubotHandler(w, r, config)
	})

	// Jenkins callback
	http.HandleFunc("/_builder", func(w http.ResponseWriter, r *http.Request) {
		builderHandler(w, r, config)
	})

	// landing page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := rootHandler(w, r, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/repositories", func(w http.ResponseWriter, r *http.Request) {
		err := repositoryHandler(w, r, config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	stats := NewStats()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, config, stats)
	})

	handler := &LoggingHandler{http.DefaultServeMux, stats}
	address := config.Address
	cert := config.CertificatePath
	key := config.KeyPath

	glog.Warningf("Starting server listening at %v.", config.Address)
	glog.Fatal(ListenAndServe(address, cert, key, handler))
}

func ListenAndServe(address, cert, key string, handler http.Handler) (err error) {
	if key != "" && cert != "" {
		err = http.ListenAndServeTLS(address, cert, key, handler)
	} else {
		err = http.ListenAndServe(address, handler)
	}

	return
}
