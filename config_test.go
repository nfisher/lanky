package main

import (
	"strings"
	"testing"
)

const validJson = `{
	"address": ":9393",
	"baseUrl": "http://ganky.local:9393/",
	"jekinsUrl": "http://jenkins.local:8080/",
	"hubot": {
		"user": "hubot",
		"password": "secret"
	},
	"github": {
		"user": "natbobc",
		"password": "secret",
		"hookSecret": "abc123"
	}
}`

const invalidJson = `{
	"address": ":9393"
	"baseUrl": "http://ganky.local:9393/",
}`

func Test_valid_json(t *testing.T) {
	c := &Config{}
	r := strings.NewReader(validJson)

	err := LoadConfig(r, c)

	if err != nil {
		t.Fatalf("err = %v, want nil", err)
	}

	expect := "abc123"
	if c.Github.HookSecret != expect {
		t.Fatalf("c.Github.HookSecret = %v, want %v", c.Github.HookSecret, expect)
	}

	expect = "secret"
	if c.Hubot.Password != expect {
		t.Fatalf("c.Hubot.Password = %v, want %v", c.Hubot.Password, expect)
	}
}

func Test_invalid_json_returns_error(t *testing.T) {
	c := &Config{}
	r := strings.NewReader(invalidJson)

	err := LoadConfig(r, c)

	if err == nil {
		t.Fatal("err == nil, want error")
	}
}

func Test_nil_config_returns_error(t *testing.T) {
	r := strings.NewReader(validJson)

	err := LoadConfig(r, nil)

	if err == nil {
		t.Fatal("err == nil, want error")
	}
}
