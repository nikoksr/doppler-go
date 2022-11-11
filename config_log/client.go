package configlog

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/config/logs APIs.
type Client struct {
	Backend doppler.Backend
	Key     string
}

// Default returns a new client based on the SDK's default backend and API key.
func Default() *Client {
	return &Client{
		Backend: doppler.GetBackend(),
		Key:     doppler.Key,
	}
}

func (c Client) get(ctx context.Context, opts *doppler.ConfigLogGetOptions) (*doppler.ConfigLog, doppler.APIResponse, error) {
	var resp doppler.ConfigLogGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs/config/logs/log",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.ConfigLog, resp.APIResponse, err
}

// Get returns a config log and its respective info.
func (c Client) Get(ctx context.Context, opts *doppler.ConfigLogGetOptions) (*doppler.ConfigLog, doppler.APIResponse, error) {
	return c.get(ctx, opts)
}

// Get returns a config log and its respective info using the default client.
func Get(ctx context.Context, opts *doppler.ConfigLogGetOptions) (*doppler.ConfigLog, doppler.APIResponse, error) {
	return Default().Get(ctx, opts)
}

func (c Client) fetchList(ctx context.Context, opts *doppler.ConfigLogListOptions) ([]*doppler.ConfigLog, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.ConfigLogListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs/config/logs",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.ConfigLogs, resp.APIResponse, err
}

// List returns a list of config logs.
func (c Client) List(ctx context.Context, opts *doppler.ConfigLogListOptions) ([]*doppler.ConfigLog, doppler.APIResponse, error) {
	return c.fetchList(ctx, opts)
}

// List returns a list of config logs using the default client.
func List(ctx context.Context, opts *doppler.ConfigLogListOptions) ([]*doppler.ConfigLog, doppler.APIResponse, error) {
	return Default().List(ctx, opts)
}

func (c Client) rollback(ctx context.Context, opts *doppler.ConfigLogRollbackOptions) (*doppler.ConfigLog, doppler.APIResponse, error) {
	var resp doppler.ConfigLogRollbackResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs/config/logs/log/rollback",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.ConfigLog, resp.APIResponse, err
}

// Rollback rolls back a config log to a previous version.
func (c Client) Rollback(ctx context.Context, opts *doppler.ConfigLogRollbackOptions) (*doppler.ConfigLog, doppler.APIResponse, error) {
	return c.rollback(ctx, opts)
}

// Rollback rolls back a config log to a previous version using the default client.
func Rollback(ctx context.Context, opts *doppler.ConfigLogRollbackOptions) (*doppler.ConfigLog, doppler.APIResponse, error) {
	return Default().Rollback(ctx, opts)
}
