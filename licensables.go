// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"errors"
	"iter"

	"github.com/thehale/lawyer/internal/paths"
)

// Licensables are the files a Declaration governs the headers of.
type Licensables struct {
	repository Repository
	paths      []Path
}

// All yields each licensable in turn, or an error for a path it couldn't
// gather or read.
func (l Licensables) All() iter.Seq2[Licensable, error] {
	return func(yield func(Licensable, error) bool) {
		files, err := l.files()
		proceed := err == nil || yield(Licensable{}, err)
		for index := 0; proceed && index < len(files); index++ {
			licensable, isLicensable, err := licensableAt(l.repository, files[index])
			switch {
			case err != nil:
				proceed = yield(Licensable{}, err)
			case isLicensable:
				proceed = yield(licensable, nil)
			}
		}
	}
}

// Check finds how each licensable's header falls short of declaration.
func (l Licensables) Check(declaration Declaration) ([]Violation, error) {
	var violations []Violation
	failures := []error{declaration.Validate()}
	for licensable, err := range l.All() {
		if err == nil {
			violations = append(violations, licensable.Violations(declaration)...)
		}
		failures = append(failures, err)
	}
	return violations, errors.Join(failures...)
}

// Fix replaces each header that falls short of declaration with the one it
// calls for.
func (l Licensables) Fix(declaration Declaration) (Changes, error) {
	var changes Changes
	failures := []error{declaration.Validate()}
	for licensable, err := range l.All() {
		if err == nil && licensable.Violations(declaration) != nil {
			var fix Changes
			fix, err = licensable.Fix(declaration)
			changes = changes.Union(fix)
		}
		failures = append(failures, err)
	}
	return changes, errors.Join(failures...)
}

func (l Licensables) files() ([]Path, error) {
	if l.paths == nil {
		return paths.FilesUnder(l.repository.root)
	} else {
		var files []Path
		var failures []error
		for _, target := range l.paths {
			under, err := filesIn(target)
			files = append(files, under...)
			failures = append(failures, err)
		}
		return files, errors.Join(failures...)
	}
}

func filesIn(path Path) ([]Path, error) {
	if path.IsDirectory() {
		return paths.FilesUnder(path)
	} else {
		return []Path{path}, nil
	}
}
