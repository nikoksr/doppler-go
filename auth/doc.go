/*
Package auth provides a client for the Doppler API's auth endpoints.

API-Docs: https://docs.doppler.com/reference/auth-revoke

Example:

	// Revoke a token.
	_, err := auth.Revoke(context.Background(), &doppler.AuthRevokeOptions{
		Tokens: []string{"token1", "token2"},
	})
	if err != nil {
		log.Fatal(err)
	}
*/
package auth
