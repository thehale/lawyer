// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages

import "github.com/thehale/lawyer/internal/syntax"

func init() {
	syntax.Register(syntax.Language{
		Name:       "CMake",
		Filenames:  []string{"CMakeLists.txt"},
		Extensions: []string{"cmake", "cmake.in"},
		Comment:    syntax.Hashes,
	})
}
