/*
Package activitylog provides a client for the Doppler API's logs endpoints.

API-Docs: https://docs.doppler.com/reference/activity-logs-list

Example:

	// Fetch activity logs
	activityLogs, _, err := activitylog.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Print activity logs
	for _, log := range activityLogs {
		fmt.Printf("%+v\n", log)
	}
*/
package activitylog
