package main_test

import "testing"

const validJson = `{
	"baseUrl": "http://ganky.local:9393/",
	"jekinsUrl": "http://jenkins.local:8080/",
	"hubot": {
		"user": "hubot",
		"password": "secret"
	}
	"github": {
		"user": "natbobc",
		"password": "secret",
		"hookSecret": "abc123"
	}
}`

func Test_fail(t *testing.T) {
	t.Fail()
}
