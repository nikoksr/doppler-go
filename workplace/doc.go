/*
Package workplace provides a client for the Doppler API's workplace endpoints.

API-Docs: https://docs.doppler.com/reference/workplace-settings-retrieve

Example:

	// Get the current workplace
	wp, _, err := workplace.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", wp)
*/
package workplace
