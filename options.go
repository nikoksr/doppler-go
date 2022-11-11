package doppler

import (
	"net/url"

	"github.com/google/go-querystring/query"
)

type (
	// ListOptions is the base struct for all list options. It contains the common parameters used across all list
	// endpoints. It's meant to be embedded in more specific list options structs.
	ListOptions struct {
		Page    int `url:"page,omitempty"`
		PerPage int `url:"per_page,omitempty"`
	}

	// parameters is a normalized way to represent query parameters. It's commonly used to represent endpoint options
	// in a more normalized way, which are then used by the Backend to make API requests. Options in contrast are a more
	// endpoint-specific way to represent them.
	parameters = url.Values
)

// extractQueryParameters transforms a struct into a url.Values object. It's commonly used to transform endpoint options
// into query parameters, which are then used by the Backend to make API requests. parameters are a normalized way to
// represent query parameters, whereas Options are a more endpoint-specific way to represent them. A nil struct will
// return an empty parameters object and no error.
func extractQueryParameters(v any) (parameters, error) {
	if v == nil {
		return make(parameters), nil
	}

	return query.Values(v)
}
