/*
Package secret provides a client for the Doppler API's secret endpoints.

API-Docs: https://docs.doppler.com/reference/config-secret-list

Example:

	  // Download secrets in docker format.
		secrets, _, err := secret.Download(context.Background(), &doppler.SecretDownloadOptions{
			Project: "your-project",
			Config:  "your-config",
			Format:  stringPointer("docker"),
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(secrets)
*/
package secret
