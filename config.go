package main

import (
	"encoding/json"
	"io"
)

type Github struct {
	User       string
	Password   string
	HookSecret string
	ApiUrl     string
}

type Hubot struct {
	User     string
	Password string
}

// Lanky run-time configuration.
type Config struct {
	Address         string
	BaseUrl         string
	CertificatePath string
	KeyPath         string
	ChatDefaultRoom string
	DatabaseUrl     string
	JenkinsUrl      string
	TemplatesDir    string
	Hubot           *Hubot
	Github          *Github
}

func LoadConfig(r io.Reader, c *Config) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(c)

	return err
}
