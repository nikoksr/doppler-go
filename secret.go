package doppler

type (
	// SecretValue represents a single Doppler secret value.
	SecretValue struct {
		Raw      *string `json:"raw,omitempty"`      // The raw value of the secret.
		Computed *string `json:"computed,omitempty"` // The computed value of the secret.
	}

	// Secret represents a single Doppler secret, including its name and value.
	Secret struct {
		Name  *string      `json:"name,omitempty"`  // The name of the secret.
		Value *SecretValue `json:"value,omitempty"` // The value of the secret.
	}

	// SecretGetResponse represents a response from the secrets get endpoint.
	//
	// Method: GET
	// Endpoint: https://api.doppler.com/v3/configs/config/secret
	// Docs:     https://docs.doppler.com/reference/config-secret-retrieve
	SecretGetResponse struct {
		APIResponse `json:",inline"`
		Secret      *Secret `json:"secret,inline"`
	}

	// SecretGetOptions represents options for the secrets get endpoint.
	SecretGetOptions struct {
		Project string `url:"project" json:"-"` // The name of the project containing the secret.
		Config  string `url:"config" json:"-"`  // The name of the config containing the secret.
		Name    string `url:"name" json:"-"`    // The name of the secret.
	}

	// SecretListResponse represents a response from the secrets list endpoint.
	//
	// Method: GET
	// Endpoint: https://api.doppler.com/v3/configs/config/secrets
	// Docs:     https://docs.doppler.com/reference/config-secret-list
	SecretListResponse struct {
		APIResponse `json:",inline"`
		Secrets     map[string]*SecretValue `json:"secrets"`
	}

	// SecretListOptions represents options for the secrets list endpoint.
	SecretListOptions struct {
		Project           string  `url:"project" json:"-"`                           // The name of the project containing the secret.
		Config            string  `url:"config" json:"-"`                            // The name of the config containing the secret.
		IncludeDynamic    *bool   `url:"include_dynamic_secrets,omitempty" json:"-"` // Whether to include dynamic secrets.
		DynamicTTLSeconds *int32  `url:"dynamic_secrets_ttl_sec,omitempty" json:"-"` // The number of seconds until dynamic leases expire. Must be used with include_dynamic_secrets.
		Secrets           *string `url:"secrets,omitempty" json:"-"`                 // A comma-separated list of secret names to include.
	}

	// SecretUpdateResponse represents a response from the secrets update endpoint.
	//
	// Method: PUT
	// Endpoint: https://api.doppler.com/v3/configs/config/secrets
	// Docs:     https://docs.doppler.com/reference/config-secret-update
	SecretUpdateResponse struct {
		APIResponse `json:",inline"`
		Secrets     map[string]string `json:"secrets"`
	}

	// SecretUpdateOptions represents options for the secrets update endpoint.
	SecretUpdateOptions struct {
		Project    string            `url:"-" json:"project"` // The name of the project containing the secret.
		Config     string            `url:"-" json:"config"`  // The name of the config containing the secret.
		NewSecrets map[string]string `url:"-" json:"secrets"` // The secrets to update.
	}

	// SecretDownloadOptions represents options for the secrets download endpoint.
	//
	// Method: GET
	// Endpoint: https://api.doppler.com/v3/configs/config/secrets/download
	// Docs:     https://docs.doppler.com/reference/config-secret-download
	SecretDownloadOptions struct {
		Project           string  `url:"project" json:"-"`                           // The name of the project containing the secret.
		Config            string  `url:"config" json:"-"`                            // The name of the config containing the secret.
		IncludeDynamic    *bool   `url:"include_dynamic_secrets,omitempty" json:"-"` // Whether to include dynamic secrets.
		DynamicTTLSeconds *int32  `url:"dynamic_secrets_ttl_sec,omitempty" json:"-"` // The number of seconds until dynamic leases expire. Must be used with include_dynamic_secrets.
		Format            *string `url:"format,omitempty" json:"-"`                  // The format to download the secrets in. See official docs for supported formats.
		NameTransformer   *string `url:"name_transformer,omitempty" json:"-"`        // The name transformer to use when downloading the secrets. See official docs for supported transformers.
	}
)
