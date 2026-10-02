// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/thehale/lawyer/internal/files"
	"github.com/thehale/lawyer/internal/spdx"
)

// A Licensing is a repository's LICENSE files.
type Licensing struct {
	repository Repository
	files      []licenseFile
	err        error
}

// A slot is where one of a declaration's licenses goes, and the files there.
type slot struct {
	license spdx.License
	path    Path
	present []licenseFile
}

const uncalledFor = "not a file the --license expression calls for"

// Check finds how the LICENSE files fall short of declaration, and any
// error reading them.
func (l Licensing) Check(declaration Declaration) ([]Violation, error) {
	expectation := declaration.licenseExpectation(l.repository)
	slots, strays := l.slots(declaration)
	var violations []Violation
	for _, slot := range slots {
		switch {
		case len(slot.present) > 1:
			for _, file := range slot.present {
				violations = append(violations, violationsAt(file.path, fmt.Sprintf("one of several files for %s", slot.license.ID))...)
			}
		case len(slot.present) == 0:
			violations = append(violations, violationsAt(slot.path, "missing")...)
		case !slot.present[0].regular:
			violations = append(violations, violationsAt(slot.present[0].path, "not a regular file")...)
		default:
			violations = append(violations, violationsAt(slot.present[0].path, slot.present[0].violations(slot.license, expectation)...)...)
		}
	}
	for _, stray := range strays {
		violations = append(violations, violationsAt(stray.path, uncalledFor)...)
	}
	return violations, errors.Join(l.err, declaration.Validate())
}

// Replace writes declared's LICENSE files in place of these, and removes those
// it doesn't hold.
func (l Licensing) Replace(declared Licensing) (Changes, error) {
	var changes Changes
	var failures []error
	for _, file := range declared.files {
		written, err := l.write(file)
		changes, failures = changes.Union(written), append(failures, err)
	}
	for _, file := range l.files {
		if file.regular && declared.fileAt(file.path) == nil {
			removed, err := remove(file)
			changes, failures = changes.Union(removed), append(failures, err)
		}
	}
	return changes, errors.Join(failures...)
}

// Fix replaces the LICENSE files with the ones declaration calls for.
func (l Licensing) Fix(declaration Declaration) (Changes, error) {
	changes, err := l.Replace(declaration.LicensingFor(l.repository))
	touched := slices.Concat(changes.Fixed, changes.Removed)
	violations, violationsErr := l.Check(declaration)
	changes.Unfixed = slices.DeleteFunc(violations, func(violation Violation) bool { return slices.Contains(touched, violation.Path) })
	return changes, errors.Join(violationsErr, err)
}

func (l Licensing) declaredBy(declaration Declaration) Licensing {
	expectation := declaration.licenseExpectation(l.repository)
	slots, _ := l.slots(declaration)
	declared := Licensing{repository: l.repository}
	for _, slot := range slots {
		switch {
		case len(slot.present) == 0:
			declared.files = append(declared.files, licenseFile{path: slot.path, content: licenseFile{}.canonical(slot.license, expectation), regular: true})
		case len(slot.present) == 1 && slot.present[0].regular && slot.present[0].violations(slot.license, expectation) != nil:
			existing := slot.present[0]
			declared.files = append(declared.files, licenseFile{name: existing.name, path: existing.path, content: existing.canonical(slot.license, expectation), regular: true})
		default:
			declared.files = append(declared.files, slot.present...)
		}
	}
	return declared
}

func (l Licensing) slots(declaration Declaration) ([]slot, []licenseFile) {
	licenses, _ := declaration.License.Licenses()
	inSlot := make([]bool, len(l.files))
	var slots []slot
	for _, license := range licenses {
		names := licenseNames(license, licenses)
		place := slot{license: license, path: l.repository.root.Child(names[0])}
		for index, file := range l.files {
			if hasName(file, names) {
				place.present, inSlot[index] = append(place.present, file), true
			}
		}
		slots = append(slots, place)
	}
	var strays []licenseFile
	for index, file := range l.files {
		if !inSlot[index] {
			strays = append(strays, file)
		}
	}
	return slots, strays
}

func (l Licensing) fileAt(path Path) *licenseFile {
	index := slices.IndexFunc(l.files, func(file licenseFile) bool { return file.path == path })

	if index >= 0 {
		return &l.files[index]
	} else {
		return nil
	}
}

func (l Licensing) write(declared licenseFile) (Changes, error) {
	existing := l.fileAt(declared.path)

	switch {
	case existing == nil:
		return fixOf(declared.path, files.Create(declared.path, declared.content))
	case existing.regular && existing.content != declared.content:
		return fixOf(declared.path, files.Replace(declared.path, existing.content, declared.content))
	default:
		return Changes{}, nil
	}
}

func remove(file licenseFile) (Changes, error) {
	err := files.Remove(file.path, file.content)

	if err == nil {
		return Changes{Removed: []Path{file.path.Printable()}}, nil
	} else {
		return Changes{}, err
	}
}

func licenseFiles(root Path) ([]licenseFile, error) {
	entries, err := os.ReadDir(string(root))
	failures := []error{err}
	var found []licenseFile
	for _, entry := range entries {
		if isManagedName(entry.Name()) {
			file, err := readLicenseFile(root, entry.Name())
			found = append(found, file)
			failures = append(failures, err)
		}
	}
	return found, errors.Join(failures...)
}
