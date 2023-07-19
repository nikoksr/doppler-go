package doppler

type (
	// ServiceToken represents a Doppler service token.
	ServiceToken struct {
		Name        *string `json:"name,omitempty"`        // Name of the service token.
		Slug        *string `json:"slug,omitempty"`        // A unique identifier of the service token.
		Key         *string `json:"key,omitempty"`         // An API key that is used for authentication. Only available when creating the token.
		Project     *string `json:"project,omitempty"`     // Unique identifier for the project object.
		Environment *string `json:"environment,omitempty"` // Unique identifier for the environment object.
		Config      *string `json:"config,omitempty"`      // The config's name.
		Access      *string `json:"access,omitempty"`      // The access level of the service token. One of read, read/write.
		ExpiresAt   *string `json:"expires_at,omitempty"`  // Date and time of the token's expiration, or null if token does not auto-expire.
		CreatedAt   *string `json:"created_at,omitempty"`  // Date and time of the object's creation.
	}

	// ServiceTokenListResponse represents a response from the service-token list endpoint.
	//
	// Method: GET
	// Endpoint: https://api.doppler.com/v3/configs/config/tokens
	// Docs:     https://docs.doppler.com/reference/config-token-list
	ServiceTokenListResponse struct {
		APIResponse `json:",inline"`
		Tokens      []*ServiceToken `json:"tokens"`
	}

	// ServiceTokenListOptions represents the options for the service-token list endpoint.
	ServiceTokenListOptions struct {
		Project string `url:"project,omitempty" json:"-"` // Unique identifier for the project object.
		Config  string `url:"config,omitempty" json:"-"`  // The config's name.
	}

	// ServiceTokenCreateResponse represents a response from the service-token create endpoint.
	//
	// Method: POST
	// Endpoint: https://api.doppler.com/v3/configs/config/tokens
	// Docs:     https://docs.doppler.com/reference/config-token-create
	ServiceTokenCreateResponse struct {
		APIResponse `json:",inline"`
		Token       *ServiceToken `json:"token,omitempty"`
	}

	// ServiceTokenCreateOptions represents the options for the service-token create endpoint.
	ServiceTokenCreateOptions struct {
		Project   string  `url:"-" json:"project,omitempty"`    // Unique identifier for the project object.
		Config    string  `url:"-" json:"config,omitempty"`     // The config's name.
		Name      string  `url:"-" json:"name,omitempty"`       // Name of the service token.
		Access    *string `url:"-" json:"access,omitempty"`     // The access level of the service token. One of read, read/write.
		ExpiresAt *string `url:"-" json:"expires_at,omitempty"` // Date and time of the token's expiration, or null if token does not auto-expire.
	}

	// ServiceTokenDeleteResponse represents a response from the service-token delete endpoint.
	//
	// Method: DELETE
	// Endpoint: https://api.doppler.com/v3/configs/config/tokens/token
	// Docs:     https://docs.doppler.com/reference/config-token-delete
	ServiceTokenDeleteResponse struct {
		APIResponse `json:",inline"`
	}

	// ServiceTokenDeleteOptions represents the options for the service-token delete endpoint.
	ServiceTokenDeleteOptions struct {
		Project string `url:"-" json:"project,omitempty"` // Unique identifier for the project object.
		Config  string `url:"-" json:"config,omitempty"`  // The config's name.
		Slug    string `url:"-" json:"slug,omitempty"`    // A unique identifier of the service token.
	}
)
