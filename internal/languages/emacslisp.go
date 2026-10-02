// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "Emacs Lisp",
		Extensions: []string{"el"},
		Comment:    syntax.Semicolons,
		Prolog:     []syntax.Prolog{syntax.Modeline},
	})
}
