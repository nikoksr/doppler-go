package secret

import (
	"context"
	"io"
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

func (c Client) get(ctx context.Context, opts *doppler.SecretGetOptions) (*doppler.Secret, doppler.APIResponse, error) {
	var resp doppler.SecretGetResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs/config/secret",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Secret, resp.APIResponse, err
}

// Get returns a config secret and its respective info.
func (c Client) Get(ctx context.Context, opts *doppler.SecretGetOptions) (*doppler.Secret, doppler.APIResponse, error) {
	return c.get(ctx, opts)
}

// Get returns a config secret and its respective info using the default client.
func Get(ctx context.Context, opts *doppler.SecretGetOptions) (*doppler.Secret, doppler.APIResponse, error) {
	return Default().Get(ctx, opts)
}

func (c Client) fetchList(ctx context.Context, opts *doppler.SecretListOptions) (map[string]*doppler.SecretValue, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.SecretListResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs/config/secrets",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Secrets, resp.APIResponse, err
}

// List returns a list of config secrets.
func (c Client) List(ctx context.Context, opts *doppler.SecretListOptions) (map[string]*doppler.SecretValue, doppler.APIResponse, error) {
	return c.fetchList(ctx, opts)
}

// List returns a list of config secrets using the default client.
func List(ctx context.Context, opts *doppler.SecretListOptions) (map[string]*doppler.SecretValue, doppler.APIResponse, error) {
	return Default().List(ctx, opts)
}

func (c Client) update(ctx context.Context, opts *doppler.SecretUpdateOptions) (map[string]string, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.SecretUpdateResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPut,
		Path:    "/v3/configs/config/secrets",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Secrets, resp.APIResponse, err
}

// Update updates a config secret.
func (c Client) Update(ctx context.Context, opts *doppler.SecretUpdateOptions) (map[string]string, doppler.APIResponse, error) {
	return c.update(ctx, opts)
}

// Update updates a config secret using the default client.
func Update(ctx context.Context, opts *doppler.SecretUpdateOptions) (map[string]string, doppler.APIResponse, error) {
	return Default().Update(ctx, opts)
}

func (c Client) download(ctx context.Context, opts *doppler.SecretDownloadOptions) (string, doppler.APIResponse, error) {
	// Make the request.
	var resp doppler.APIResponse
	httpResp, err := c.Backend.CallRaw(ctx, &doppler.Request{
		Method:  http.MethodGet,
		Path:    "/v3/configs/config/secrets/download",
		Key:     c.Key,
		Payload: opts,
	})
	if err != nil {
		return "", resp, err
	}
	defer httpResp.Body.Close()

	// Fill the APIResponse object with details from the HTTP response.
	resp.WithDetails(httpResp)

	// Check if the APIResponse object indicates an error.
	if resp.Error() != nil {
		return "", resp, resp.Error()
	}

	// Read the response body. This can be of different formats depending on the request. The format may bet
	// set in the doppler.SecretDownloadOptions.Format field. Hence, we're treating it as a byte array.
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return "", resp, err
	}

	return string(body), resp, nil
}

// Download downloads a config secret.
func (c Client) Download(ctx context.Context, opts *doppler.SecretDownloadOptions) (string, doppler.APIResponse, error) {
	return c.download(ctx, opts)
}

// Download downloads a config secret using the default client.
func Download(ctx context.Context, opts *doppler.SecretDownloadOptions) (string, doppler.APIResponse, error) {
	return Default().Download(ctx, opts)
}
