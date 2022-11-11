/*
Package share provides a client for the Doppler API's share endpoints.

API-Docs: https://docs.doppler.com/reference/share-secret

Example:

		// Share a plain text secret.
		sharedSecret, _, err := share.PlainSecret(context.Background(), &doppler.SharePlainOptions{
	    Secret: "my-secret",
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(sharedSecret)
*/
package share
