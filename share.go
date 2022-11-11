package doppler

type (
	// SharePlain allows to generate a Doppler Share link by sending a plain text secret.
	SharePlain struct {
		URL              *string `json:"url,omitempty"`
		AuthenticatedURL *string `json:"authenticated_url,omitempty"`
		Password         *string `json:"password,omitempty"`
	}

	// ShareEncrypted allows to generate a Doppler Share link by sending an end-to-end encrypted secret.
	ShareEncrypted struct {
		URL *string `json:"url,omitempty"`
	}

	// SharePlainResponse represents a response from the share plain endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v1/share/secrets/plain
	// Docs:     https://docs.doppler.com/reference/share-secret
	SharePlainResponse struct {
		APIResponse `json:",inline"`
		Secret      *SharePlain `json:",inline"`
	}

	// SharePlainOptions represents the options for the share plain endpoint.
	SharePlainOptions struct {
		Secret      string `url:"-" json:"secret" validate:"required"` // Plain text secret to share.
		ExpireViews *int32 `url:"-" json:"expire_views,omitempty"`     // Number of views before the link expires. Valid ranges: 1 to 50. -1 for unlimited.
		ExpireDays  *int32 `url:"-" json:"expire_days,omitempty"`      // Number of days before the link expires. Valid range: 1 to 90.
	}

	// ShareEncryptedResponse represents a response from the share encrypted endpoint.
	//
	// Method:   POST
	// Endpoint: https://api.doppler.com/v1/share/secrets/encrypted
	// Docs:     https://docs.doppler.com/reference/share-secret-encrypted
	ShareEncryptedResponse struct {
		APIResponse `json:",inline"`
		Secret      *ShareEncrypted `json:",inline"`
	}

	// ShareEncryptedOptions represents the options for the share encrypted endpoint.
	ShareEncryptedOptions struct {
		Secret      string `url:"-" json:"encrypted_secret" validate:"required"`                 // Base64 encoded AES-GCM encrypted secret to share. See docs for more details.
		Password    string `url:"-" json:"hashed_password" validate:"required"`                  // SHA256 hash of the password. This is NOT the hash of the derived encryption key.
		KDF         string `url:"-" json:"encryption_kdf" validate:"required,eq=pbkdf2"`         // The key derivation function used. Must by "pbkdf2".
		SaltRounds  int32  `url:"-" json:"encryption_salt_rounds" validate:"required,eq=100000"` // Number of salt rounds used by KDF. Must be "100000".
		ExpireViews *int32 `url:"-" json:"expire_views,omitempty"`                               // Number of views before the link expires. Valid ranges: 1 to 50. -1 for unlimited.
		ExpireDays  *int32 `url:"-" json:"expire_days,omitempty"`                                // Number of days before the link expires. Valid range: 1 to 90.
	}
)

const (
	// EncryptionKDF is the key derivation function used for encrypted secrets. As stated in the docs [1] this
	// has to be "pbkdf2". This is a constatnt to avoid typos and help with testing.
	//
	// [1]: https://docs.doppler.com/reference/share-secret-encrypted
	EncryptionKDF = "pbkdf2"

	// EncryptionSaltRounds is the number of salt rounds used by the key derivation function. As stated in the docs [1] this
	// has to be "100000". This is a constatnt to avoid typos and help with testing.
	//
	// [1]: https://docs.doppler.com/reference/share-secret-encrypted
	EncryptionSaltRounds = 100000
)
