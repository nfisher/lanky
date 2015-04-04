package main

import (
	"encoding/json"
	"io"
	"time"
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

type Jenkins struct {
	BaseUrl  string
	TrayFeed string
}

// Lanky run-time configuration.
type Config struct {
	Address         string
	BaseUrl         string
	CertificatePath string
	KeyPath         string
	ChatDefaultRoom string
	DatabaseUrl     string
	TemplatesDir    string
	Jenkins         *Jenkins
	Hubot           *Hubot
	Github          *Github
}

func (c *Config) ClientTimeout() time.Duration {
	return time.Duration(5 * time.Second)
}

func (c *Config) TrayFeedUrl() string {
	if c.Jenkins == nil {
		return "http://localhost:8080"
	}

	return c.Jenkins.BaseUrl + c.Jenkins.TrayFeed
}

func LoadConfig(r io.Reader, c *Config) error {
	dec := json.NewDecoder(r)
	err := dec.Decode(c)

	return err
}
