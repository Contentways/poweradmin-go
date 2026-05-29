// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import "net/http"

// Response wraps [http.Response] with parsed metadata extracted from the
// Poweradmin API envelope (currently: pagination).
type Response struct {
	*http.Response
	Meta ResponseMeta
}

// ResponseMeta holds parsed metadata that accompanied an API response.
type ResponseMeta struct {
	Pagination *Pagination
}

// Pagination is the domain-level representation of pagination metadata.
type Pagination struct {
	Page     int
	PerPage  int
	Total    int
	LastPage int
}
