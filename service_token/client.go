package servicetoken

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/configs/config/secrets APIs.
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

func (c Client) fetchList(ctx context.Context, opts *doppler.ServiceTokenListOptions) ([]*doppler.ServiceToken, doppler.APIResponse, error) {
	var resp doppler.ServiceTokenListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs/config/tokens",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Tokens, resp.APIResponse, err
}

// List returns a list of service tokens.
func (c Client) List(ctx context.Context, opts *doppler.ServiceTokenListOptions) ([]*doppler.ServiceToken, doppler.APIResponse, error) {
	return c.fetchList(ctx, opts)
}

// List returns a list of config secrets using the default client.
func List(ctx context.Context, opts *doppler.ServiceTokenListOptions) ([]*doppler.ServiceToken, doppler.APIResponse, error) {
	return Default().List(ctx, opts)
}

func (c Client) create(ctx context.Context, opts *doppler.ServiceTokenCreateOptions) (*doppler.ServiceToken, doppler.APIResponse, error) {
	var resp doppler.ServiceTokenCreateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs/config/tokens",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Token, resp.APIResponse, err
}

// Create creates a new service tokens.
func (c Client) Create(ctx context.Context, opts *doppler.ServiceTokenCreateOptions) (*doppler.ServiceToken, doppler.APIResponse, error) {
	return c.create(ctx, opts)
}

// Create creates a new service tokens using the default client.
func Create(ctx context.Context, opts *doppler.ServiceTokenCreateOptions) (*doppler.ServiceToken, doppler.APIResponse, error) {
	return Default().Create(ctx, opts)
}

func (c Client) delete(ctx context.Context, _ *doppler.ServiceTokenDeleteOptions) (doppler.APIResponse, error) {
	var resp doppler.ServiceTokenDeleteResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method: http.MethodDelete,
		Path:   "/v3/configs/config/tokens/token",
		Key:    c.Key,
	}, &resp)

	return resp.APIResponse, err
}

// Delete deletes a service tokens.
func (c Client) Delete(ctx context.Context, opts *doppler.ServiceTokenDeleteOptions) (doppler.APIResponse, error) {
	return c.delete(ctx, opts)
}

// Delete deletes a service tokens using the default client.
func Delete(ctx context.Context, opts *doppler.ServiceTokenDeleteOptions) (doppler.APIResponse, error) {
	return Default().Delete(ctx, opts)
}
