package client

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"infopuller/internal/utils/config"
)

var (
	ErrNewRequest = fmt.Errorf("failed to build a request")
	ErrDo         = fmt.Errorf("request failed")
	ErrReadAll    = fmt.Errorf("failed to read a body")
)

type Client struct {
	UnimplementedClient

	http.Client

	Log *slog.Logger

	Config config.Config
}

func New(log *slog.Logger, c config.Config) *Client {
	const op = "client.New()"

	log.Debug(
		"initiallizing client",
		slog.String("op", op),
	)

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
	const op = "client.Shutdown()"

	c.Log.Debug(
		"stopping the client",
		slog.String("op", op),
	)

	c.CloseIdleConnections()
}

func (c *Client) Random() ([]byte, error) {
	const op = "client.Random()"

	c.Log.Debug(
		"request received by client random handler",
		slog.String("op", op),
	)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, c.Config.Client.RandomURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNewRequest, err)
	}

	req = addHeaders(req, c.Config.Client.KinopoiskAPIKey)

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDo, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %v", ErrDo, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrReadAll, err)
	}

	return body, nil
}

func addHeaders(req *http.Request, key string) *http.Request {
	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-KEY", key)

	return req
}

type UnimplementedClient struct{}

func (u *UnimplementedClient) Random() ([]byte, error) {
	return nil, nil
}
