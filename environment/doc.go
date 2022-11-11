/*
Package environments provides a client for the Doppler API's environments endpoints.

API-Docs: https://docs.doppler.com/reference/environment-object

Example:

	// Fetch all environments of a project
	environments, _, err := environment.List(ctx, &doppler.EnvironmentListOptions{
		Project: "my-project",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print the slug of each environment
	for _, env := range environments {
		fmt.Println(env.ID)
	}
*/
package environment
