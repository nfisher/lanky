package main

type Config struct {
	BaseUrl         string
	ChatDefaultRoom string
	DatabaseUrl     string
	JenkinsUrl      string
	TemplatesDir    string

	Hubot struct {
		User     string
		Password string
	}

	Github struct {
		User       string
		Password   string
		HookSecret string
		ApiUrl     string
	}
}
