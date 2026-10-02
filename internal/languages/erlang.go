// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "Erlang",
		Extensions: []string{"erl"},
		Shebangs:   []string{"escript"},
		Comment:    syntax.Percents,
		Prolog:     []syntax.Prolog{syntax.Shebang},
	})
}
