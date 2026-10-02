// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"errors"
	"fmt"

	"github.com/thehale/lawyer/internal/history"
	"github.com/thehale/lawyer/internal/paths"
)

// A Path is the path of a file or directory.
type Path = paths.Path

// A Repository is the working tree whose files, LICENSE files and git history
// lawyer reads.
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

// Licensing is the repository's LICENSE files.
func (r Repository) Licensing() Licensing {
	files, err := licenseFiles(r.root)
	return Licensing{repository: r, files: files, err: err}
}

// Check finds how every licensable's header and the repository's Licensing
// fall short of declaration.
func (r Repository) Check(declaration Declaration) ([]Violation, error) {
	err := declaration.Validate()

	if err == nil {
		headers, headersErr := r.Licensables().Check(declaration)
		licenses, licensesErr := r.Licensing().Check(declaration)
		return append(headers, licenses...), errors.Join(headersErr, licensesErr)
	} else {
		return nil, err
	}
}

// Fix makes every licensable's header and the repository's Licensing what
// declaration calls for.
func (r Repository) Fix(declaration Declaration) (Changes, error) {
	err := declaration.Validate()

	if err == nil {
		headers, headersErr := r.Licensables().Fix(declaration)
		licenses, licensesErr := r.Licensing().Fix(declaration)
		return headers.Union(licenses), errors.Join(headersErr, licensesErr)
	} else {
		return Changes{}, err
	}
}

// Warnings are what keeps the repository from checking declaration fully.
func (r Repository) Warnings(declaration Declaration) []string {
	if declaration.Years == "" && !r.history.IsEveryEditDated() {
		return []string{fmt.Sprintf("%s, so years are only checked for being in the past", r.history.Gap)}
	} else {
		return nil
	}
}
