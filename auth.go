package doppler

import "encoding/json"

type (
	// Auth is the object representing an auth.
	Auth struct{}

	// AuthToken is the object representing an auth token.
	AuthToken struct {
		Token *string `url:"-" json:"token,omitempty"` // The token itself.
	}

	// AuthRevokeResponse is the response from the AuthRevokeOptions endpoint.
	//
	// Method:    POST
	// Endpoint:  https://api.doppler.com/v3/auth/revoke
	// Docs:      https://docs.doppler.com/reference/auth-revoke
	AuthRevokeResponse struct {
		APIResponse `json:",inline"`
	}

	// AuthRevokeOptions revokes an auth token.
	AuthRevokeOptions struct {
		Tokens []AuthToken `url:"-" json:"tokens" validate:"gt=0"` // A list of tokens to revoke.
	}
)

// MarshalJSON is a custom JSON marshaller for AuthRevokeOptions. The API expects the Tokens slice directly, instead of
// wrapped in a Tokens object. To not break the API and compatibility with this library, we need to do this custom
// marshalling.
func (opts *AuthRevokeOptions) MarshalJSON() ([]byte, error) {
	return json.Marshal(opts.Tokens)
}
