/*
Package servicetoken provides a client for the Doppler API's service token endpoints.

API-Docs: https://docs.doppler.com/reference/config-token-list

Example:

	  // List all service tokens
		tokens, _, err := servicetoken.List(context.Background(), &doppler.ServiceTokenListOptions{
			Project: "your-project",
			Config:  "your-config",
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(tokens)
*/
package servicetoken
