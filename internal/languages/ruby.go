// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "Ruby",
		Filenames:  []string{"Berksfile", "Brewfile", "Capfile", "Dangerfile", "Fastfile", "Gemfile", "Guardfile", "Podfile", "Rakefile", "Steepfile", "Thorfile", "Vagrantfile"},
		Extensions: []string{"rb", "gemspec", "rake", "ru"},
		Shebangs:   []string{"ruby"},
		Comment:    syntax.Hashes,
		Prolog:     []syntax.Prolog{syntax.Shebang, syntax.Rack, syntax.Encoding, syntax.MagicComments},
	})
}
