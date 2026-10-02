// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"fmt"

	"github.com/thehale/lawyer/internal/history"
	"github.com/thehale/lawyer/internal/paths"
)

// A Path is the path of a file or directory.
type Path = paths.Path

// A Repository is the working tree whose files and git history lawyer reads.
type Repository struct {
	root    Path
	history history.History
}

// NewRepository opens the repository at root, inside git or not.
func NewRepository(root Path) Repository {
	return Repository{root: root, history: history.Of(root)}
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

// Warnings are what keeps the repository from checking declaration fully.
func (r Repository) Warnings(declaration Declaration) []string {
	if declaration.Years == "" && !r.history.IsEveryEditDated() {
		return []string{fmt.Sprintf("%s, so years are only checked for being in the past", r.history.Gap)}
	} else {
		return nil
	}
}
