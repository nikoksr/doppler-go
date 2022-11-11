/*
Package doppler provides a Go client library for the Doppler REST API.

Website:             https://www.doppler.com/
Docs:                https://docs.doppler.com/docs
API Reference:       https://docs.doppler.com/reference/api
Auth Token Formats:  https://docs.doppler.com/reference/auth-token-formats

Example:

	    package main

	    import (
	        "github.com/nikoksr/doppler-go"
		    "github.com/nikoksr/doppler-go/project"
	    )

			func main() {
				ctx := context.Background()

				// Set your API key
				doppler.Key = "YOUR_API_KEY"

				// Get a list of all projects
				projects, _,  err := project.List(ctx, nil)
				if err != nil {
					panic(err)
				}

				// Print the names of all projects and update them afterwards
				for _, project := range projects {
					fmt.Println(project.Name)
				}

				// Update the first project
				_, _, err = project.Update(ctx, &doppler.ProjectUpdateOptions{
					Name: projects[0].Name,
					NewConfig:        ...,              // Leaving this out, so nobody accidentally overwrites a real project
					NewDescription:   ...,
				})
				if err != nil {
					panic(err)
				}
			}
*/
package doppler
