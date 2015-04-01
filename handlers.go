package main

import (
	"html/template"
	"net/http"
	"time"
)

const rootHtml = `<!DOCTYPE html>
<html lang="en">
  <head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
	<link href="//fonts.googleapis.com/css?family=Raleway:400,300,600" rel="stylesheet" type="text/css">
	<style>
	html {
		font-size:62.5%;
	}
	body {
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
		background:#eee;
		display:block;
		color:#1EAEDB;
		text-decoration:none;
		text-indent:1rem;
	}
	a:hover {
		background:#ccc;
	}
	</style>
	</head>
	<body>
	<h1>Lanky</h1>
	<ul>
	{{range .Project}}
	<li><a href="{{.WebUrl}}{{.LastBuildLabel}}/console">{{.BuildTime}} - {{.Name}}</a> 
	{{end}}
	</ul>
	</body>
</html>`

var rootTemplate = template.Must(template.New("root").Parse(rootHtml))

func rootHandler(w http.ResponseWriter, r *http.Request, config *Config) {
	client := http.Client{
		// TODO: Should probably make this configurable via the config.
		Timeout: time.Duration(5 * time.Second),
	}

	resp, err := client.Get(config.JenkinsUrl + "/cc.xml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p := &Projects{}
	defer resp.Body.Close()
	err = ReadTrayFeed(resp.Body, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = rootTemplate.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
