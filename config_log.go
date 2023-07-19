package doppler

type (
	// ConfigLogDiff represents a diff between two config logs.
	ConfigLogDiff struct {
		Name  *string `json:"name,omitempty"`  // Name of the config.
		Added *string `json:"added,omitempty"` // Added data.
	}

	// ConfigLog represents a Doppler config log.
	ConfigLog struct {
		ID          *string         `json:"id,omitempty"`          // Unique identifier of the config log.
		Text        *string         `json:"text,omitempty"`        // Text describing the event.
		HTML        *string         `json:"html,omitempty"`        // HTML describing the event.
		Diff        []ConfigLogDiff `json:"diff,omitempty"`        // Diff between the previous and current config.
		Rollback    *bool           `json:"rollback,omitempty"`    // Is this config log a rollback of a previous log.
		User        *User           `json:"user,omitempty"`        // User that triggered the event.
		Project     *string         `json:"project,omitempty"`     // Identifier of the project that the config belongs to.
		Environment *string         `json:"environment,omitempty"` // Identifier of the environment that the config belongs to.
		Config      *string         `json:"config,omitempty"`      // Name of the config.
		CreatedAt   *string         `json:"created_at,omitempty"`  // Date and time of the object's creation.
	}

	// ConfigLogGetResponse represents a response from the config log get endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/configs/config/logs/log
	// Docs:     https://docs.doppler.com/reference/config-log-retrieve
	ConfigLogGetResponse struct {
		APIResponse `json:",inline"`
		ConfigLog   *ConfigLog `json:"config_log,omitempty"`
	}

	// ConfigLogGetOptions represents the options for the config log get endpoint.
	ConfigLogGetOptions struct {
		Project string `url:"project" json:"-"` // Identifier of the project that the config belongs to.
		Config  string `url:"config" json:"-"`  // Name of the config.
		ID      string `url:"log" json:"-"`     // Unique identifier of the config log.
	}

	// ConfigLogListResponse represents a response from the config log list endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/configs/config/logs
	// Docs:     https://docs.doppler.com/reference/config-log-list
	ConfigLogListResponse struct {
		APIResponse `json:",inline"`
		ConfigLogs  []*ConfigLog `json:"logs"`
	}

	// ConfigLogListOptions represents the query parameters for a config log list request.
	ConfigLogListOptions struct {
		ListOptions `url:",inline" json:"-"`
		Project     string `url:"project" json:"-"` // Identifier of the project that the config belongs to.
		Config      string `url:"config" json:"-"`  // Name of the config.
	}

	// ConfigLogRollbackResponse represents a response from the config log rollback endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/configs/config/logs/log/rollback
	// Docs:     https://docs.doppler.com/reference/config-log-rollback
	ConfigLogRollbackResponse struct {
		APIResponse `json:",inline"`
		ConfigLog   *ConfigLog `json:"config_log,omitempty"`
	}

	// ConfigLogRollbackOptions represents the options for the config log rollback endpoint.
	ConfigLogRollbackOptions struct {
		Project string `url:"project" json:"-"` // Identifier of the project that the config belongs to.
		Config  string `url:"config" json:"-"`  // Name of the config.
		ID      string `url:"log" json:"-"`     // Unique identifier of the config log.
	}
)
