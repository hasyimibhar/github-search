package github

import (
	"net"
	"net/http"
	"time"

	"github.com/hasyimibhar/github-search/common"
)

const (
	BaseURL string = "https://api.github.com"
)

// ClientConfig configures the Github API client.
type ClientConfig struct {
	ClientID     string
	ClientSecret string
}

// Client is the Github API client.
type Client struct {
	HTTPClient *http.Client
	Logger     common.Logger

	config *ClientConfig
}

// NewClient creates a new client.
//
// If config is nil, requests are not authenticated.
func NewClient(config *ClientConfig, logger common.Logger) *Client {
	ts := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   true,
	}

	return &Client{
		HTTPClient: &http.Client{
			Timeout:   time.Second * 30,
			Transport: ts,
		},
		Logger: logger,

		config: config,
	}
}

// Search returns the search API gateway.
func (c *Client) Search() SearchGateway {
	return SearchGateway{HTTPClient: c.HTTPClient, config: c.config, logger: c.Logger}
}

// Error is the error returned by the client.
type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
