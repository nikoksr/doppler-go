/*
Package audit provides a client for the Doppler API's audit endpoints.

API-Docs: https://docs.doppler.com/reference/audit-api-reference

Example:

	    // Fetch audit logs for the current workplace
			logs, _, err := audit.WorkplaceGet(context.Background(), &doppler.AuditWorkplaceGetOptions{})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(logs)
*/
package audit
