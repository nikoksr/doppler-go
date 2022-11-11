/*
Package project provides a client for the Doppler API's project endpoints.

API-Docs: https://docs.doppler.com/reference/project-object

Example:

	// Fetch a list of projects.
	list, _, err := project.List(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print the list of projects.
	fmt.Printf("%+v", list)
*/
package project
