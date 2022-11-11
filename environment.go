package doppler

type (
	// Environment represents a doppler environment.
	Environment struct {
		ID             *string `json:"id,omitempty"`               // An identifier for the object.
		Slug           *string `json:"slug,omitempty"`             // A unique identifier for the environment.
		Name           *string `json:"name,omitempty"`             // Name of the environment.
		Project        *string `json:"project,omitempty"`          // Identifier of the project the environment belongs to.
		InitialFetchAt *string `json:"initial_fetch_at,omitempty"` // Date and time of the first secrets fetch from a config in the environment.
		CreatedAt      *string `json:"created_at,omitempty"`       // Date and time of the object's creation.
	}

	// EnvironmentGetResponse represents a response from the environment get endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/environments/environment
	// Docs:     https://docs.doppler.com/reference/environment-retrieve
	EnvironmentGetResponse struct {
		APIResponse `json:",inline"`
		Environment *Environment `json:"environment,omitempty"`
	}

	// EnvironmentGetOptions represents the options for the environment get endpoint.
	EnvironmentGetOptions struct {
		Project string `url:"project" json:"-" validate:"required"`     // Identifier of the project the environment belongs to.
		Slug    string `url:"environment" json:"-" validate:"required"` // A unique identifier for the environment.
	}

	// EnvironmentListResponse represents a response from the environment list endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/environments
	// Docs:     https://docs.doppler.com/reference/environment-list
	EnvironmentListResponse struct {
		APIResponse  `json:",inline"`
		Environments []*Environment `json:"environments"`
	}

	// EnvironmentListOptions represents the query parameters for a environment list request.
	EnvironmentListOptions struct {
		Project string `url:"project" json:"-" validate:"required"` // Identifier of the project the environment belongs to.
	}

	// EnvironmentCreateResponse represents a response from the environment create endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/environments
	// Docs:     https://docs.doppler.com/reference/environment-create
	EnvironmentCreateResponse struct {
		APIResponse `json:",inline"`
		Environment *Environment `json:"environment"`
	}

	// EnvironmentCreateOptions represents the body parameters for a environment create request.
	EnvironmentCreateOptions struct {
		Project string `url:"project" json:"-" validate:"required"`        // Identifier of the project the environment belongs to.
		Name    string `url:"-" json:"name,omitempty" validate:"required"` // Name of the environment.
		Slug    string `url:"-" json:"slug,omitempty" validate:"required"` // A unique identifier for the environment.
	}

	// EnvironmentRenameResponse represents a doppler environment rename request.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/environments/environment
	// Docs:     https://docs.doppler.com/reference/environment-rename
	EnvironmentRenameResponse struct {
		APIResponse `json:",inline"`
		Environment *Environment `json:"environment"`
	}

	// EnvironmentRenameOptions represents the body parameters for a environment rename request.
	EnvironmentRenameOptions struct {
		Project string  `url:"project" json:"-" validate:"required"`     // Identifier of the project the environment belongs to.
		Slug    string  `url:"environment" json:"-" validate:"required"` // A unique identifier for the environment.
		NewName *string `url:"-" json:"name,omitempty"`                  // New name of the environment.
		NewSlug *string `url:"-" json:"slug,omitempty"`                  // New slug of the environment.
	}

	// EnvironmentDeleteResponse represents a response from the environment delete endpoint.
	//
	// Method:   DELETE
	// Endpoint: https://api.doppler.com/v3/environments/environment
	// Docs:     https://docs.doppler.com/reference/environment-delete
	EnvironmentDeleteResponse struct {
		APIResponse `json:",inline"`
	}

	// EnvironmentDeleteOptions represents the body parameters for a environment delete request.
	EnvironmentDeleteOptions struct {
		Project string `url:"project" json:"-" validate:"required"`     // Identifier of the project the environment belongs to.
		Slug    string `url:"environment" json:"-" validate:"required"` // A unique identifier for the environment.
	}
)
