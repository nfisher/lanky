package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

const defaultConfigPath = "lanky.json"

type RuntimeStats struct {
	Started   time.Time
	status1xx uint64
	status2xx uint64
	status3xx uint64
	status4xx uint64
	status5xx uint64
}

func (rs *RuntimeStats) StartDate() string {
	return rs.Started.Format("2006-01-02 15:04")
}

func (rs *RuntimeStats) inc(addr *uint64) {
	atomic.AddUint64(addr, 1)
}

func (rs *RuntimeStats) IncStatus(code int) {
	codeClass := code / 100
	switch codeClass {
	case 1:
		rs.Inc1xx()
		return
	case 2:
		rs.Inc2xx()
		return
	case 3:
		rs.Inc3xx()
		return
	case 4:
		rs.Inc4xx()
		return
	case 5:
		rs.Inc5xx()
		return
	}

	log.Fatalf("Unexpected response code %v.", code)
}

func (rs *RuntimeStats) Inc1xx() { rs.inc(&rs.status1xx) }
func (rs *RuntimeStats) Inc2xx() { rs.inc(&rs.status2xx) }
func (rs *RuntimeStats) Inc3xx() { rs.inc(&rs.status3xx) }
func (rs *RuntimeStats) Inc4xx() { rs.inc(&rs.status4xx) }
func (rs *RuntimeStats) Inc5xx() { rs.inc(&rs.status5xx) }

func (rs *RuntimeStats) Status1xx() uint64 { return atomic.LoadUint64(&rs.status1xx) }
func (rs *RuntimeStats) Status2xx() uint64 { return atomic.LoadUint64(&rs.status2xx) }
func (rs *RuntimeStats) Status3xx() uint64 { return atomic.LoadUint64(&rs.status3xx) }
func (rs *RuntimeStats) Status4xx() uint64 { return atomic.LoadUint64(&rs.status4xx) }
func (rs *RuntimeStats) Status5xx() uint64 { return atomic.LoadUint64(&rs.status5xx) }

func NewStats() *RuntimeStats {
	return &RuntimeStats{
		Started: time.Now(),
	}
}

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

	Serve(config)
}

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
		lh.RuntimeStats.IncStatus(bw.Status)
	}
	log.Printf("%v %v - \"%v %v %v\" %v %v", ip, r.Header.Get("User-Agent"), r.Method, r.URL.Path, r.Proto, bw.Status, bw.Wrote)
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

	stats := NewStats()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, config, stats)
	})

	handler := &LoggingHandler{http.DefaultServeMux, stats}
	address := config.Address
	cert := config.CertificatePath
	key := config.KeyPath

	log.Printf("Starting server listening at %v.", config.Address)
	log.Fatal(ListenAndServe(address, cert, key, handler))
}

func ListenAndServe(address, cert, key string, handler http.Handler) (err error) {
	if key != "" && cert != "" {
		err = http.ListenAndServeTLS(address, cert, key, handler)
	} else {
		err = http.ListenAndServe(address, handler)
	}

	return
}
