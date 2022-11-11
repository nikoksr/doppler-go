package dynamicsecret

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /configs/config/dynamic_secrets/dynamic_secret APIs.
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

func (c Client) issueLease(ctx context.Context, opts *doppler.DynamicSecretIssueLeaseOptions) (doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.DynamicSecretIssueLeaseResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v3/configs/config/dynamic_secrets/dynamic_secret/leases",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.APIResponse, err
}

// IssueLease issues a lease for a dynamic secret.
func (c Client) IssueLease(ctx context.Context, opts *doppler.DynamicSecretIssueLeaseOptions) (doppler.APIResponse, error) {
	return c.issueLease(ctx, opts)
}

// IssueLease issues a lease for a dynamic secret using the default client.
func IssueLease(ctx context.Context, opts *doppler.DynamicSecretIssueLeaseOptions) (doppler.APIResponse, error) {
	return Default().IssueLease(ctx, opts)
}

func (c Client) revokeLease(ctx context.Context, opts *doppler.DynamicSecretRevokeLeaseOptions) (doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.DynamicSecretRevokeLeaseResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodDelete,
		Path:    "/v3/configs/config/dynamic_secrets/dynamic_secret/leases/lease",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.APIResponse, err
}

// RevokeLease revokes a lease for a dynamic secret.
func (c Client) RevokeLease(ctx context.Context, opts *doppler.DynamicSecretRevokeLeaseOptions) (doppler.APIResponse, error) {
	return c.revokeLease(ctx, opts)
}

// RevokeLease revokes a lease for a dynamic secret using the default client.
func RevokeLease(ctx context.Context, opts *doppler.DynamicSecretRevokeLeaseOptions) (doppler.APIResponse, error) {
	return Default().RevokeLease(ctx, opts)
}
