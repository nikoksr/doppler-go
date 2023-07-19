package doppler

type (
	// DynamicSecretIssueLeaseResponse represents a response from the dynamic secret issue lease endpoint.
	//
	// Method: POST
	// Endpoint: https://api.doppler.com/v3/configs/config/dynamic_secrets/dynamic_secret/leases
	// Docs: https://docs.doppler.com/reference/dynamic-secret-issue-lease
	DynamicSecretIssueLeaseResponse struct {
		APIResponse `json:",inline"`
	}

	// DynamicSecretIssueLeaseOptions represents the options for the dynamic secret issue lease endpoint.
	DynamicSecretIssueLeaseOptions struct {
		Project    string `url:"-" json:"project"`        // The project where the dynamic secret is located
		Config     string `url:"-" json:"config"`         // The config where the dynamic secret is located
		Name       string `url:"-" json:"dynamic_secret"` // The dynamic secret to issue a lease for
		TTLSeconds int32  `url:"-" json:"ttl_seconds"`    // The number of seconds the lease should last
	}

	// DynamicSecretRevokeLeaseResponse represents a response from the dynamic secret revoke lease endpoint.
	//
	// Method: DELETE
	// Endpoint:  https://api.doppler.com/v3/configs/config/dynamic_secrets/dynamic_secret/leases/lease
	// Docs: https://docs.doppler.com/reference/dynamic-secret-issue-revoke-lease
	DynamicSecretRevokeLeaseResponse struct {
		APIResponse `json:",inline"`
	}

	// DynamicSecretRevokeLeaseOptions represents the options for the dynamic secret revoke lease endpoint.
	DynamicSecretRevokeLeaseOptions struct {
		Project string `url:"-" json:"project"`        // The project where the dynamic secret is located
		Config  string `url:"-" json:"config"`         // The config where the dynamic secret is located
		Name    string `url:"-" json:"dynamic_secret"` // The dynamic secret to revoke a lease for
		Slug    string `url:"-" json:"slug"`           // The lease to revoke
	}
)
