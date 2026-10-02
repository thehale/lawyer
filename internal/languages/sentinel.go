// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "Sentinel",
		Extensions: []string{"sentinel"},
		Comment:    syntax.Mixed,
		Prolog:     []syntax.Prolog{syntax.Policy},
	})
}
