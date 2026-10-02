// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	_ "github.com/thehale/lawyer/internal/languages"
	"github.com/thehale/lawyer/internal/spdx"
	"github.com/thehale/lawyer/internal/syntax"
)

// A Language is a kind of file lawyer recognizes, and what it recognizes it by.
type Language struct {
	Name string
	// Filenames are whole file names, such as "Dockerfile".
	Filenames []string
	// Extensions come without the dot, such as "py" or "cmake.in".
	Extensions []string
	// Shebangs are the interpreters that select an extensionless script, as
	// named after #! or #!/usr/bin/env, such as "python3".
	Shebangs []string
	// HasHeader is false for a language with no comment syntax, such as
	// JSON, whose files lawyer recognizes and leaves alone.
	HasHeader bool
}

// A LicenseID is an SPDX short identifier, such as "MPL-2.0".
type LicenseID = spdx.ID

// Languages lists every Language lawyer recognizes, ordered by name.
func Languages() []Language {
	var languages []Language
	for _, language := range syntax.KnownLanguages() {
		languages = append(languages, Language{
			Name:       language.Name,
			Filenames:  language.Filenames,
			Extensions: language.Extensions,
			Shebangs:   language.Shebangs,
			HasHeader:  language.HasHeader(),
		})
	}
	return languages
}

// Licenses lists the license ids an Expression may use: every license SPDX
// marks OSI-approved and not deprecated.
func Licenses() []LicenseID {
	return spdx.IDs()
}
