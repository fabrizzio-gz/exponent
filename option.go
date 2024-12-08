package exponent

import "net/http"

type Config struct {
	Host       string
	ApiURL     string
	AcessToken string
	HttpClient *http.Client
}

type Option func(*Config)

func WithHost(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

func WithApiURL(apiURL string) Option {
	return func(c *Config) {
		c.ApiURL = apiURL
	}
}

func WithAccessToken(accessToken string) Option {
	return func(c *Config) {
		c.AcessToken = accessToken
	}
}

func WithHttpClient(httpClient *http.Client) Option {
	return func(c *Config) {
		c.HttpClient = httpClient
	}
}

func withDefaults(c *Config) {
	if c.Host == "" {
		c.Host = "https://exp.host"
	}
	if c.ApiURL == "" {
		c.ApiURL = "/--/api/v2"
	}
	if c.HttpClient == nil {
		c.HttpClient = &http.Client{}
	}
}
