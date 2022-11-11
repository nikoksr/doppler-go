/*
Package config provides a client for the Doppler API's configs endpoints.

API-Docs: https://docs.doppler.com/reference/config-object

Example:

	// Fetch a list of configs for a project.
	configs, _, err := config.List(context.Background(), &doppler.ConfigListOptions{
		Project: "my-project",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print the configs
	for _, config := range configs {
		fmt.Printf("Config: %s", config.Name)
	}
*/
package config
