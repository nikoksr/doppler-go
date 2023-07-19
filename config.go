package doppler

type (
	// Config represents a Doppler configuration.
	Config struct {
		Name           *string `json:"name,omitempty"`             // Name of the configuration.
		Project        *string `json:"project,omitempty"`          // Identifier of the project that the config belongs to.
		Environment    *string `json:"environment,omitempty"`      // Identifier of the environment that the config belongs to.
		Root           *bool   `json:"root,omitempty"`             // Whether the config is the root of the environment.
		Locked         *bool   `json:"locked,omitempty"`           // Whether the config can be renamed and/or deleted.
		InitialFetchAt *string `json:"initial_fetch_at,omitempty"` // Date and time of the first secrets fetch.
		LastFetchAt    *string `json:"last_fetch_at,omitempty"`    // Date and time of the last secrets fetch.
		CreatedAt      *string `json:"created_at,omitempty"`       // Date and time of the object's creation.
	}

	// ConfigGetResponse represents a response from the config get endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/configs/config
	// Docs:     https://docs.doppler.com/reference/config-retrieve
	ConfigGetResponse struct {
		APIResponse `json:",inline"`
		Config      *Config `json:"config,omitempty"`
	}

	// ConfigGetOptions represents the options for the config get endpoint.
	ConfigGetOptions struct {
		Project string `url:"project" json:"-"` // Identifier of the project that the config belongs to.
		Config  string `url:"config" json:"-"`  // Name of the config.
	}

	// ConfigListResponse represents a response from the config list endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/configs
	// Docs:     https://docs.doppler.com/reference/config-list
	ConfigListResponse struct {
		APIResponse `json:",inline"`
		Configs     []*Config `json:"configs"`
	}

	// ConfigListOptions represents the query parameters for a config list request.
	ConfigListOptions struct {
		ListOptions `url:",inline" json:"-"`
		Project     string `url:"project" json:"-"` // Identifier of the project that the config belongs to.
	}

	// ConfigCreateResponse represents a response from the config create endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/configs
	// Docs:     https://docs.doppler.com/reference/config-create
	ConfigCreateResponse struct {
		APIResponse `json:",inline"`
		Config      *Config `json:"config"`
	}

	// ConfigCreateOptions represents the body parameters for a config create request.
	ConfigCreateOptions struct {
		Project     string `url:"-" json:"project"`     // Identifier of the project that the config belongs to.
		Environment string `url:"-" json:"environment"` // Identifier of the environment that the config belongs to.
		Name        string `url:"-" json:"name"`        // Name of the new branch configuration.
	}

	// ConfigUpdateResponse represents a doppler config update request.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/configs/config
	// Docs:     https://docs.doppler.com/reference/config-update
	ConfigUpdateResponse struct {
		APIResponse `json:",inline"`
		Config      *Config `json:"config"`
	}

	// ConfigUpdateOptions represents the body parameters for a config update request.
	ConfigUpdateOptions struct {
		Project string `url:"-" json:"project"` // Identifier of the project that the config belongs to.
		Config  string `url:"-" json:"config"`  // Name of the config.
		NewName string `url:"-" json:"name"`    // New name of the config.
	}

	// ConfigDeleteResponse represents a response from the config delete endpoint.
	//
	// Method:   DELETE
	// Endpoint: https://api.doppler.com/v3/configs/config
	// Docs:     https://docs.doppler.com/reference/config-delete
	ConfigDeleteResponse struct {
		APIResponse `json:",inline"`
	}

	// ConfigDeleteOptions represents the body parameters for a config delete request.
	ConfigDeleteOptions struct {
		Project string `url:"-" json:"project"` // Identifier of the project that the config belongs to.
		Config  string `url:"-" json:"config"`  // Name of the config.
	}

	// ConfigLockResponse represents a response from the config lock endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/configs/config/lock
	// Docs:     https://docs.doppler.com/reference/config-lock
	ConfigLockResponse struct {
		APIResponse `json:",inline"`
		Config      *Config `json:"config"`
	}

	// ConfigLockOptions represents the body parameters for a config lock request.
	ConfigLockOptions struct {
		Project string `url:"-" json:"project"` // Identifier of the project that the config belongs to.
		Config  string `url:"-" json:"config"`  // Name of the config.
	}

	// ConfigUnlockResponse represents a response from the config unlock endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/configs/config/unlock
	// Docs:     https://docs.doppler.com/reference/config-unlock
	ConfigUnlockResponse struct {
		APIResponse `json:",inline"`
		Config      *Config `json:"config"`
	}

	// ConfigUnlockOptions represents the body parameters for a config unlock request.
	ConfigUnlockOptions struct {
		Project string `url:"-" json:"project"` // Identifier of the project that the config belongs to.
		Config  string `url:"-" json:"config"`  // Name of the config.
	}

	// ConfigCloneResponse represents a response from the config clone endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/configs/config/clone
	// Docs:     https://docs.doppler.com/reference/config-clone
	ConfigCloneResponse struct {
		APIResponse `json:",inline"`
		Config      *Config `json:"config"`
	}

	// ConfigCloneOptions represents the body parameters for a config clone request.
	ConfigCloneOptions struct {
		Project   string `url:"-" json:"project"` // Identifier of the project that the config belongs to.
		Config    string `url:"-" json:"config"`  // Name of the config.
		NewConfig string `url:"-" json:"name"`    // Name of the new config.
	}
)
