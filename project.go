package doppler

type (
	// Project represents a doppler project.
	Project struct {
		ID          *string `json:"id,omitempty"`          // ID is the unique identifier for the object.
		Name        *string `json:"name,omitempty"`        // Name is the name of the project.
		Slug        *string `json:"slug,omitempty"`        // Slug is an abbreviated name for the project.
		Description *string `json:"description,omitempty"` // Description is the description of the project.
		CreatedAt   *string `json:"created_at,omitempty"`  // CreatedAt is the time the project was created.
	}

	// ProjectGetResponse represents a response from the project get endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/projects/project
	// Docs:     https://docs.doppler.com/reference/project-retrieve
	ProjectGetResponse struct {
		APIResponse `json:",inline"`
		Project     *Project `json:"project,omitempty"`
	}

	// ProjectGetOptions represents the options for the project get endpoint.
	ProjectGetOptions struct {
		Name string `url:"project" json:"-" validate:"required"` // Name is the name of the project.
	}

	// ProjectListResponse represents a response from the project list endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/projects
	// Docs:     https://docs.doppler.com/reference/project-list
	ProjectListResponse struct {
		APIResponse `json:",inline"`
		Projects    []*Project `json:"projects"`
	}

	// ProjectListOptions represents the query parameters for a project list request.
	ProjectListOptions struct {
		ListOptions `url:",inline" json:"-"`
	}

	// ProjectCreateResponse represents a response from the project create endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/projects
	// Docs:     https://docs.doppler.com/reference/project-create
	ProjectCreateResponse struct {
		APIResponse `json:",inline"`
		Project     *Project `json:"project"`
	}

	// ProjectCreateOptions represents the body parameters for a project create request.
	ProjectCreateOptions struct {
		Name        string  `url:"-" json:"name" validate:"required"` // Name of the project.
		Description *string `url:"-" json:"description,omitempty"`    // Description of the project.
	}

	// ProjectUpdateResponse represents a doppler project update request.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/projects/project
	// Docs:     https://docs.doppler.com/reference/project-update
	ProjectUpdateResponse struct {
		APIResponse `json:",inline"`
		Project     *Project `json:"project"`
	}

	// ProjectUpdateOptions represents the body parameters for a project update request.
	ProjectUpdateOptions struct {
		Name           string  `url:"-" json:"project" validate:"required"`        // Name of the project.
		NewName        string  `url:"-" json:"name,omitempty" validate:"required"` // New name of the project.
		NewDescription *string `url:"-" json:"description,omitempty"`              // New description of the project.
	}

	// ProjectDeleteResponse represents a response from the project delete endpoint.
	//
	// Method:   DELETE
	// Endpoint: https://api.doppler.com/v3/projects/project
	// Docs:     https://docs.doppler.com/reference/project-delete
	ProjectDeleteResponse struct {
		APIResponse `json:",inline"`
	}

	// ProjectDeleteOptions represents the body parameters for a project delete request.
	ProjectDeleteOptions struct {
		Name string `url:"-" json:"project" validate:"required"` // Name of the project.
	}
)
