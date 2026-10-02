// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "Dart",
		Extensions: []string{"dart"},
		Shebangs:   []string{"dart"},
		Comment:    syntax.Slashes,
		Prolog:     []syntax.Prolog{syntax.Shebang},
	})
}
