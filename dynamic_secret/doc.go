/*
Package configlog provides a client for the Doppler API's config logs endpoints.

API-Docs: https://docs.doppler.com/reference/config-log-list

Example:

	// Fetch config logs
	configLogs, _, err := configlog.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print config logs
	for _, log := range configLogs {
		fmt.Printf("%+v\n", log)
	}
*/

package dynamicsecret
