// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"net/url"
	"strconv"
)

// ListOpts controls pagination for List operations.
// A zero value requests the API's default page and size.
type ListOpts struct {
	Page    int
	PerPage int
}

// values renders the options as query parameters.
func (o ListOpts) values() url.Values {
	v := url.Values{}
	if o.Page > 0 {
		v.Set("page", strconv.Itoa(o.Page))
	}
	if o.PerPage > 0 {
		v.Set("per_page", strconv.Itoa(o.PerPage))
	}
	return v
}

// appendQuery returns path with the encoded query appended (if any).
func appendQuery(path string, v url.Values) string {
	q := v.Encode()
	if q == "" {
		return path
	}
	if containsRune(path, '?') {
		return path + "&" + q
	}
	return path + "?" + q
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
