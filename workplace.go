package doppler

type (
	// Workplace represents a doppler workplace.
	Workplace struct {
		ID           *string `json:"id,omitempty"`            // ID of the workplace.
		Name         *string `json:"name,omitempty"`          // Name of the workplace.
		BillingEmail *string `json:"billing_email,omitempty"` // BillingEmail doppler will send invoices to.
	}

	// WorkplaceGetResponse represents a response from the workplace get endpoint.
	//
	// Method:   GET
	// Endpoint: https://api.doppler.com/v3/workplace
	// Docs:     https://docs.doppler.com/reference/workplace-settings-retrieve
	WorkplaceGetResponse struct {
		APIResponse `json:",inline"`
		Workplace   *Workplace `json:"workplace,omitempty"`
	}

	// WorkplaceUpdateResponse represents a response from the workplace update endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v3/workplace
	// Docs:     https://docs.doppler.com/reference/workplace-settings-update
	WorkplaceUpdateResponse struct {
		APIResponse `json:",inline"`
		Workplace   *Workplace `json:"workplace"`
	}

	// WorkplaceUpdateOptions represents a request to the workplace update endpoint.
	WorkplaceUpdateOptions struct {
		NewName         *string `json:"name,omitempty"`          // New name of the workplace.
		NewBillingEmail *string `json:"billing_email,omitempty"` // New billing email Doppler will send invoices to.
	}
)
