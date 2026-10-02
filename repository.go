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

// Licensables are the files at paths. They're read as they're iterated.
func (r Repository) Licensables(paths ...Path) Licensables {
	return Licensables{repository: r, paths: paths}
}
