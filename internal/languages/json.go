// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

// JSON has no comment syntax, so it takes no header.
func init() {
	syntax.Register(syntax.Language{
		Name:       "JSON",
		Extensions: []string{"json"},
	})
}
