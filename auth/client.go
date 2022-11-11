package auth

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v3/auth APIs.
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

func (c Client) revoke(ctx context.Context, opts *doppler.AuthRevokeOptions) (doppler.APIResponse, error) {
	var resp doppler.AuthRevokeResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/auth/revoke",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.APIResponse, err
}

// Revoke revokes auth tokens.
func (c Client) Revoke(ctx context.Context, opts *doppler.AuthRevokeOptions) (doppler.APIResponse, error) {
	return c.revoke(ctx, opts)
}

// Revoke revokes auth tokens using the default client.
func Revoke(ctx context.Context, opts *doppler.AuthRevokeOptions) (doppler.APIResponse, error) {
	return Default().Revoke(ctx, opts)
}
