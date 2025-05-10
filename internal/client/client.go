package client

import (
	"context"
	"fmt"
	"infopuller/internal/utils/config"
	"io"
	"log/slog"
	"net/http"
)

// TODO: FIGURE OUT BETTER ERRORS

type Client struct {
	http.Client

	Log *slog.Logger

	Config config.Config
}

func New(log *slog.Logger, c config.Config) *Client {
	client := http.Client{
		Timeout: c.Client.Timeout,
	}

	// TODO: DEBUG LOG CLIENT START

	return &Client{
		Client: client,

		Log: log,

		Config: c,
	}
}

func (c *Client) Shutdown() {
	c.CloseIdleConnections()
}

func (c *Client) Random() ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, c.Config.Client.RandomURL, nil)
	if err != nil {
		// TODO: LOG ERROR

		return nil, err
	}

	req = addHeaders(req, c.Config.Client.KinopoiskAPIKey)

	resp, err := c.Do(req)
	if err != nil {
		// TODO: LOG ERROR

		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		// TODO: LOG ERROR

		return nil, fmt.Errorf("request failed: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func addHeaders(req *http.Request, key string) *http.Request {
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-KEY", key)

	return req
}
