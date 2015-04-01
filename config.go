package main

import (
	"encoding/json"
	"io"
)

// Lanky run-time configuration.
type Config struct {
	Address         string
	BaseUrl         string
	ChatDefaultRoom string
	DatabaseUrl     string
	JenkinsUrl      string
	TemplatesDir    string

	Hubot *struct {
		User     string
		Password string
	}

	Github *struct {
		User       string
		Password   string
		HookSecret string
		ApiUrl     string
	}
}

func LoadConfig(r io.Reader, c *Config) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(c)

	return err
}
