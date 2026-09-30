// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"errors"
	"iter"
	"slices"

	"github.com/thehale/lawyer/internal/paths"
)

// Licensables are the files a Declaration governs the headers of.
type Licensables struct {
	repository Repository
	paths      []Path
	globs      []Glob
}

// A Glob matches paths, with * within a directory and ** across them, as in
// "**/vendor/**".
type Glob = paths.Glob

// Excluding leaves out the files any of the globs match.
func (l Licensables) Excluding(globs ...Glob) Licensables {
	l.globs = append(slices.Clone(l.globs), globs...)
	return l
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
	files, err := l.gather()
	return slices.DeleteFunc(files, paths.Matcher(l.globs)), err
}

func (l Licensables) gather() ([]Path, error) {
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
		within, outside := paths.FilesWithin(l.repository.root, files)
		return paths.Distinct(within), errors.Join(append(failures, outside)...)
	}
}

func filesIn(path Path) ([]Path, error) {
	switch {
	case path == "":
		return nil, nil
	case path.IsDirectory():
		return paths.FilesUnder(path)
	default:
		return []Path{path}, nil
	}
}
