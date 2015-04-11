package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const githubEventType = "X-GitHub-Event"
const githubSignature = "X-Hub-Signature"
const githubSignaturePrefix = "sha1="
const githubUserAgent = "GitHub-Hookshot/"
const statusHtml = `<!DOCTYPE html>
<html lang="en">
  <head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Lanky</title>
	<style>
	html {
		font-size:62.5%;
	}
	body {
		color:#222;
		font-family: HelveticaNeue, 'Helvetica Neue', Helvetica, Arial, sans-serif;
		font-size:1.5em;
		margin:1rem auto;
		position:relative;
		width:960px;
	}
	ul {
		margin:0;
		padding:0;
	}
	li {
		list-style:none;
		line-height:4rem;
		height:4rem;
		margin-bottom:1px;
	}
	th {
		text-align:left;
	}
	.number {
		text-align:right;
	}
	.key {
		min-width:10rem;
	}
	</style>
	</head>
	<body>
	<h1>Lanky</h1>
	<table>
	<tr><th class=key>Key</th><th>Value</th></tr>
	<tr><td>Started</td><td>{{.StartDate}}</td></tr>
	<tr><td>Version</td><td class=number>{{.Version}}</td></tr>
	<tr><td># Goroutines</td><td class=number>{{.NumGoroutine}}</td></tr>
	<tr><td>1XX</td><td class=number>{{.Status1xx}}</td></tr>
	<tr><td>2XX</td><td class=number>{{.Status2xx}}</td></tr>
	<tr><td>3XX</td><td class=number>{{.Status3xx}}</td></tr>
	<tr><td>4XX</td><td class=number>{{.Status4xx}}</td></tr>
	<tr><td>5XX</td><td class=number>{{.Status5xx}}</td></tr>
	<tr><td>Bytes from System</td><td class=number>{{.Sys}}</td></tr>
	<tr><td>Heap in Use</td><td class=number>{{.HeapInuse}}</td></tr>
	<tr><td>Heap System</td><td class=number>{{.HeapSys}}</td></tr>
	<tr><td>Total Allocation</td><td class=number>{{.TotalAlloc}}</td></tr>
	</table>
	</body>
</html>`

const rootHtml = `<!DOCTYPE html>
<html lang="en">
  <head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
	<title>Lanky</title>
	<link href="//fonts.googleapis.com/css?family=Raleway:400,300,600" rel="stylesheet" type="text/css">
	<style>
	html {
		font-size:62.5%;
	}
	body {
		color:#222;
		font-family: Raleway, HelveticaNeue, 'Helvetica Neue', Helvetica, Arial, sans-serif;
		font-size:1.5em;
		margin:1rem auto;
		position:relative;
		width:960px;
	}
	ul {
		margin:0;
		padding:0;
	}
	li {
		list-style:none;
		line-height:4rem;
		height:4rem;
		margin-bottom:1px;
	}
	a {
		color:#1EAEDB;
	}
	ul a {
		background:#eee;
		display:block;
		text-decoration:none;
		text-indent:1rem;
	}
	.Success a {
		background:#517F1A;
		color:white;
	}
	.Failure a {
		background:#B2123F;
		color:white;
	}
	a:hover {
		background:#ccc;
	}
	</style>
	</head>
	<body>
	<h1>Lanky</h1>
	<p>Last {{.Len}} builds sorted by:
	{{if .ByDate}}
	date, <a href="?by=status">status</a>
	{{else}}
	<a href="?by=date">date</a>, status
	{{end}}
	</p>
	<ul>
	{{range .Project}}
	<li class="{{.LastBuildStatus}}"><a href="{{.ConsoleUrl}}">{{.BuildTime}} - {{.Name}} (#{{.LastBuildLabel}})</a>
	{{end}}
	</ul>
	</body>
</html>`

const repositoryHtml = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Lanky</title>
<link href="//fonts.googleapis.com/css?family=Raleway:400,300,600" rel="stylesheet" type="text/css">
<style>
html {
	font-size:62.5%;
}
body {
	color:#222;
	font-family: Raleway, HelveticaNeue, 'Helvetica Neue', Helvetica, Arial, sans-serif;
	font-size:1.5em;
	margin:1rem auto;
	position:relative;
	width:960px;
}
ul {
	margin:0;
	padding:0;
}
li {
	list-style:none;
	line-height:4rem;
	height:4rem;
	margin-bottom:1px;
}
a {
	color:#1EAEDB;
}
ul a {
	background:#eee;
	display:block;
	text-decoration:none;
	text-indent:1rem;
}
</style>
</head>
<body>
<h1>Lanky</h1>
<p>{{.Len}} repositories.</p>
<ul>
{{range .}}
<li>{{.FullName}}
{{end}}
</ul>
</body>
</html>`

var rootTemplate = template.Must(template.New("root").Parse(rootHtml))
var statusTemplate = template.Must(template.New("status").Parse(statusHtml))
var repositoryTemplate = template.Must(template.New("repository").Parse(repositoryHtml))

func statusHandler(w http.ResponseWriter, r *http.Request, config *Config, stats *RuntimeStats) error {
	stats.Update()

	// lock all of the reads to ensure a consistent point in time measurement
	stats.RLock()
	err := statusTemplate.Execute(w, stats)
	stats.RUnlock()

	if err != nil {
		return err
	}

	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request, config *Config) error {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return nil
	}

	j := NewJenkins(config)
	if j == nil {
		return errors.New("Jenkins configuration is invalid.")
	}

	by := r.URL.Query().Get("by")
	p := &Projects{}
	err := j.TrayFeed(p, by)
	if err != nil {
		return err
	}

	err = rootTemplate.Execute(w, p)
	if err != nil {
		return err
	}

	return nil
}

func builderHandler(w http.ResponseWriter, r *http.Request, config *Config) (err error) {
	http.Error(w, "Not implemented yet", http.StatusInternalServerError)
	return
}

func sign(payload []byte, key string) []byte {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write(payload)
	return mac.Sum(nil)
}

func githubHandler(w http.ResponseWriter, r *http.Request, config *Config) (err error) {
	if r.Method != "POST" {
		http.Error(w, "Unauthorized", http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.UserAgent(), githubUserAgent) {
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}

	signatureHeader := r.Header.Get(githubSignature)
	if !strings.HasPrefix(signatureHeader, githubSignaturePrefix) {
		http.Error(w, "Invalid signature.", http.StatusBadRequest)
		return
	}

	reqSignature, err := hex.DecodeString(signatureHeader[len(githubSignaturePrefix):])
	if err != nil {
		http.Error(w, "Invalid signature.", http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		http.Error(w, "Unable to read message body.", http.StatusInternalServerError)
		return
	}

	localSignature := sign(body, config.Github.HookSecret)
	if !hmac.Equal(localSignature, reqSignature) {
		http.Error(w, "Invalid signature.", http.StatusBadRequest)
		return
	}

	event := r.Header.Get(githubEventType)
	switch event {
	case "push":
		http.Error(w, "Not implemented yet", http.StatusInternalServerError)
		return
	case "ping":
		fmt.Fprint(w, "OK: 1")
		return
	}

	http.Error(w, "Invalid event type specified.", http.StatusBadRequest)
	return
}

func hubotHandler(w http.ResponseWriter, r *http.Request, config *Config) (err error) {
	http.Error(w, "Not implemented yet", http.StatusInternalServerError)
	return
}

var repos *Repositories = new(Repositories)
var lastUpdated time.Time
var reposSync sync.Mutex
var reposSwap sync.RWMutex

func repositoryHandler(w http.ResponseWriter, r *http.Request, config *Config) (err error) {
	if r.Method != "GET" {
		http.Error(w, "Unauthorized", http.StatusMethodNotAllowed)
		return
	}

	cl := NewGithub(config)
	if cl == nil {
		return errors.New("Github configuration is invalid.")
	}

	if r.URL.Query().Get("update") == "now" {
		reposSync.Lock()
		defer reposSync.Unlock()
		now := time.Now()
		isAfter := now.After(lastUpdated.Add(5 * time.Minute))

		reps := make(Repositories, 0, 100)
		if isAfter {
			err = cl.ListRepositories(config.Github.Organization, &reps)
			if err != nil {
				return err
			}

			lastUpdated = time.Now()

			reposSwap.Lock()
			repos = &reps
			reposSwap.Unlock()
		}
	}

	reposSwap.RLock()
	err = repositoryTemplate.Execute(w, *repos)
	reposSwap.RUnlock()
	if err != nil {
		return err
	}

	return nil
}
