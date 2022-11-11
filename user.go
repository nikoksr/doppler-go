package doppler

type (
	// User represents a user in Doppler.
	User struct {
		Email           *string `json:"email,omitempty"`             // The user's email address.
		Name            *string `json:"name,omitempty"`              // The user's name.
		UserName        *string `json:"username,omitempty"`          // The user's username.
		ProfileImageURL *string `json:"profile_image_url,omitempty"` // The user's profile image URL.
	}
)
