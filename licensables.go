// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"errors"
	"iter"
)

// Licensables are the files a Declaration governs the headers of.
type Licensables struct {
	repository Repository
	paths      []Path
}

// All yields each licensable in turn, or an error for a path it couldn't
// read.
func (l Licensables) All() iter.Seq2[Licensable, error] {
	return func(yield func(Licensable, error) bool) {
		proceed := true
		for index := 0; proceed && index < len(l.paths); index++ {
			licensable, isLicensable, err := licensableAt(l.repository, l.paths[index])
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
