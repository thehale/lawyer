// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "PHP",
		Extensions: []string{"php"},
		Shebangs:   []string{"php"},
		Comment:    syntax.Stars,
		Prolog:     []syntax.Prolog{syntax.Shebang, syntax.OpeningTag},
	})
}
