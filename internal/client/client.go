package client

import (
	"infopuller/internal/utils/config"
	"log/slog"
	"net/http"
)

type Client struct {
	Client http.Client

	Log *slog.Logger

	Config config.Config
}

func New(log *slog.Logger, c config.Config) *Client {
	client := http.Client{
		Timeout: c.Client.Timeout,
	}

	return &Client{
		Client: client,

		Log: log,

		Config: c,
	}
}

func (c *Client) Shutdown() {
	c.Client.CloseIdleConnections()
}
