package doppler

type (
	// AuditWorkplace represents an audit entry for a workplace
	AuditWorkplace struct {
		ID           *string `json:"id,omitempty"`            // The ID of the workplace
		Name         *string `json:"name,omitempty"`          // The name of the workplace
		BillingEmail *string `json:"billing_email,omitempty"` // The billing email of the workplace
		SAMLEnabled  *bool   `json:"saml_enabled,omitempty"`  // Whether SAML is enabled for the workplace
		SCIMEnabled  *bool   `json:"scim_enabled,omitempty"`  // Whether SCIM is enabled for the workplace
	}

	// AuditWorkplaceUser represents an audit entry for a workplace user
	AuditWorkplaceUser struct {
		ID        *string `json:"id,omitempty"`         // The ID of the user
		Access    *string `json:"access,omitempty"`     // The access level of the user
		User      *User   `json:"user,omitempty"`       // The user
		CreatedAt *string `json:"created_at,omitempty"` // The time the user was added to the workplace
	}

	// AuditWorkplaceGetResponse represents a response from the audit workplace get endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/workplace
	// Docs:     https://docs.doppler.com/reference/audit-workplace-retrieve
	AuditWorkplaceGetResponse struct {
		APIResponse    `json:",inline"`
		AuditWorkplace *AuditWorkplace `json:"workplace,omitempty"`
	}

	// AuditWorkplaceGetOptions represents options for the audit workplace get endpoint.
	AuditWorkplaceGetOptions struct {
		Settings *bool `url:"settings,omitempty" json:"-"` // If true, the api will return more information if the workplace has e.g. SAML enabled and SCIM enabled.
	}

	// AuditWorkplaceUserGetResponse represents a response from the audit workplace user get endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/workplace/users/{workplace_user_id}
	// Docs:     https://docs.doppler.com/reference/audit-workplace-user-retrieve
	AuditWorkplaceUserGetResponse struct {
		APIResponse        `json:",inline"`
		AuditWorkplaceUser *AuditWorkplaceUser `json:"workplace_user,omitempty"`
	}

	// AuditWorkplaceUserGetOptions represents options for the audit workplace user get endpoint.
	AuditWorkplaceUserGetOptions struct {
		UserID   string `url:"-" json:"-"`
		Settings *bool  `url:"settings,omitempty" json:"-"` // If true, the api will return more information if the workplace has e.g. SAML enabled and SCIM enabled.
	}

	// AuditWorkplaceUserListResponse represents a response from the audit workplace user list endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/workplace/users
	// Docs:     https://docs.doppler.com/reference/audit-workplace-users-retrieve
	AuditWorkplaceUserListResponse struct {
		APIResponse         `json:",inline"`
		AuditWorkplaceUsers []*AuditWorkplaceUser `json:"workplace_users,omitempty"`
	}

	// AuditWorkplaceUserListOptions represents options for the audit workplace user list endpoint.
	AuditWorkplaceUserListOptions struct {
		Settings *bool `url:"settings,omitempty" json:"-"` // If true, the api will return more information if the workplace has e.g. SAML enabled and SCIM enabled.
	}
)
