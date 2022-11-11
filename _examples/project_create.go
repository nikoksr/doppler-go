package main

import (
	"context"
	"log"
	"os"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/pointer"
	"github.com/nikoksr/doppler-go/project"
)

func main() {
	// Set the API key.
	doppler.Key = os.Getenv("DOPPLER_API_KEY")

	// Create a new project.
	prj, _, err := project.Create(context.Background(), &doppler.ProjectCreateOptions{
		Name:        "your-project-name", // Required; change this to your project name.
		Description: pointer.To("Specify an optional description here"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print the project's ID.
	log.Printf("Successfully created project with ID %q\n", *prj.ID)
}
