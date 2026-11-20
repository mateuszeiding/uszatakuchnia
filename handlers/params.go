// Package handlers holds the API endpoint logic. Each handler serves one method of one route;
// routes, methods and 404/405 live in router/routes.go, and CORS plus the write gate (resp.Begin)
// run once in the Vercel function in api/index.go, before the router. Handlers must NOT call
// resp.Begin (duplicate CORS headers); a test in router/ enforces it.
package handlers

import (
	"net/http"
	"strconv"
)

// pathID returns {id} from the path. The router lets only numbers ≥ 1 through (404 otherwise).
func pathID(r *http.Request) uint {
	id, _ := strconv.ParseUint(r.PathValue("id"), 10, 64)
	return uint(id)
}
