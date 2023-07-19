package doppler

type (
	// ActivityLog represents a doppler activity log.
	ActivityLog struct {
		ID          *string `json:"id,omitempty"`          // ID is the unique identifier for the activity log.
		Text        *string `json:"text,omitempty"`        // Text describing the event.
		HTML        *string `json:"html,omitempty"`        // HTML describing the event.
		User        *User   `json:"user,omitempty"`        // User is the user that triggered the event.
		Project     *string `json:"project,omitempty"`     // Project is the project that triggered the event.
		Environment *string `json:"environment,omitempty"` // Environment is the environment's unique identifier.
		Config      *string `json:"config,omitempty"`      // Config is the config's name.
		CreatedAt   *string `json:"created_at,omitempty"`  // CreatedAt is the time the activity log was created.
	}

	// ActivityLogGetResponse represents a response from the activity log endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/logs
	// Docs:     https://docs.doppler.com/reference/activity-log-retrieve
	ActivityLogGetResponse struct {
		APIResponse `json:",inline"`
		ActivityLog *ActivityLog `json:"log"`
	}

	// ActivityLogGetOptions represents the query parameters for an activity log get request.
	ActivityLogGetOptions struct {
		ID string `url:"log" json:"-"` // ID is the unique identifier for the log object.
	}

	// ActivityLogListResponse represents a response from the activity log list endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/logs
	// Docs:     https://docs.doppler.com/reference/activity-logs-list
	// ActivityLogsListResponse represents a response from the activity logs list endpoint.
	ActivityLogListResponse struct {
		APIResponse  `json:",inline"`
		ActivityLogs []*ActivityLog `json:"logs"`
	}

	// ActivityLogListOptions represents the query parameters for an activity log list request.
	ActivityLogListOptions struct {
		ListOptions `url:",inline" json:"-"`
	}
)
