// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "XML",
		Extensions: []string{"xml", "wxi", "wxl", "wxs"},
		Comment:    syntax.Markup,
		Prolog:     []syntax.Prolog{syntax.Declaration},
	})
}
