// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "JavaScript",
		Extensions: []string{"js", "mjs", "cjs", "jsx", "gjs"},
		Shebangs:   []string{"node"},
		Comment:    syntax.Stars,
		Prolog:     []syntax.Prolog{syntax.Shebang, syntax.Environment},
	})
}
