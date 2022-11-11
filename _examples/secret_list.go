package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nikoksr/doppler-go"
	"github.com/nikoksr/doppler-go/secret"
)

func main() {
	// Set the API key.
	doppler.Key = os.Getenv("DOPPLER_API_KEY")

	// Fetch all secrets for a given project's config.
	secrets, _, err := secret.List(context.Background(), &doppler.SecretListOptions{
		Project: "your-project-name", // The project name.
		Config:  "your-config-name",  // The config name for the project.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print all secrets.
	for name, value := range secrets {
		fmt.Printf("%s: %v\n", name, value)
	}
}
