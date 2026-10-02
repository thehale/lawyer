// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "PowerShell",
		Extensions: []string{"ps1", "psd1", "psm1"},
		Shebangs:   []string{"pwsh"},
		Comment:    syntax.Hashes,
		Prolog:     []syntax.Prolog{syntax.Shebang},
	})
}
