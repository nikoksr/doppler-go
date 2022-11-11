package share

import (
	"context"
	"net/http"

	"github.com/nikoksr/doppler-go"
)

// Client is the client used to invoke /v1/share/secrets/plain APIs.
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

func (c Client) plainSecret(ctx context.Context, opts *doppler.SharePlainOptions) (*doppler.SharePlain, doppler.APIResponse, error) {
	var resp doppler.SharePlainResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v1/share/secrets/plain",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Secret, resp.APIResponse, err
}

// PlainSecret generates a Doppler Share link by sending a plain text secret.
//
// Note:
func (c Client) PlainSecret(ctx context.Context, opts *doppler.SharePlainOptions) (*doppler.SharePlain, doppler.APIResponse, error) {
	return c.plainSecret(ctx, opts)
}

// PlainSecret generates a Doppler Share link by sending a plain text secret using the default client.
func PlainSecret(ctx context.Context, opts *doppler.SharePlainOptions) (*doppler.SharePlain, doppler.APIResponse, error) {
	return Default().PlainSecret(ctx, opts)
}

func (c Client) encryptedSecret(ctx context.Context, opts *doppler.ShareEncryptedOptions) (*doppler.ShareEncrypted, doppler.APIResponse, error) {
	var resp doppler.ShareEncryptedResponse
	err := c.Backend.Call(ctx, &doppler.Request{
		Method:  http.MethodPost,
		Path:    "/v1/share/secrets/encrypted",
		Key:     c.Key,
		Payload: opts,
	}, &resp)

	return resp.Secret, resp.APIResponse, err
}

// EncryptedSecret generates a Doppler Share link by sending an encrypted secret.
//
// Note: This endpoint requires you to take extra steps to ensure the security of your secret. Please follow
// the instructions in the documentation to ensure your secret is encrypted properly.
//
// Docs: https://docs.doppler.com/reference/share-secret-encrypted
func (c Client) EncryptedSecret(ctx context.Context, opts *doppler.ShareEncryptedOptions) (*doppler.ShareEncrypted, doppler.APIResponse, error) {
	return c.encryptedSecret(ctx, opts)
}

// EncryptedSecret generates a Doppler Share link by sending an encrypted secret using the default client.
//
// Note: This endpoint requires you to take extra steps to ensure the security of your secret. Please follow
// the instructions in the documentation to ensure your secret is encrypted properly.
//
// Docs: https://docs.doppler.com/reference/share-secret-encrypted
func EncryptedSecret(ctx context.Context, opts *doppler.ShareEncryptedOptions) (*doppler.ShareEncrypted, doppler.APIResponse, error) {
	return Default().EncryptedSecret(ctx, opts)
}
