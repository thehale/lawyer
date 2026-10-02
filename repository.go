// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import "github.com/thehale/lawyer/internal/paths"

// A Path is the path of a file or directory.
type Path = paths.Path

// A Repository is the working tree whose files lawyer reads.
type Repository struct {
	root Path
}

// NewRepository opens the repository at root.
func NewRepository(root Path) Repository {
	return Repository{root: root}
}

// Licensables are the files under paths, walking directories, or every file
// git tracks or doesn't ignore when no paths are given. They're read as
// they're iterated.
func (r Repository) Licensables(paths ...Path) Licensables {
	return Licensables{repository: r, paths: paths}
}

// Check finds how every licensable's header falls short of declaration.
func (r Repository) Check(declaration Declaration) ([]Violation, error) {
	return r.Licensables().Check(declaration)
}

// Fix makes every licensable's header what declaration calls for.
func (r Repository) Fix(declaration Declaration) (Changes, error) {
	return r.Licensables().Fix(declaration)
}
