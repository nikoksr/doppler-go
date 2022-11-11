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

	// Fetching secrets but ignoring them; we're only interested in the rate limit.
	_, resp, err := secret.List(context.Background(), &doppler.SecretListOptions{
		Project: "your-project-name", // The project name.
		Config:  "your-config-name",  // The config name for the project.
	})
	if err != nil {
		log.Fatal(err)
	}

	// Print the rate limit.
	fmt.Printf("Rate limit: %d/%d", resp.RateLimit.Remaining, resp.RateLimit.Limit)
	fmt.Printf("Rate limit reset: %s", resp.RateLimit.Reset)
}
