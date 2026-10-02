// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

// A Violation is one way a file falls short of a Declaration.
type Violation struct {
	Path    Path
	Message string
}

// Changes are what a fix did, and what it couldn't do: the files it Fixed and
// Removed, and the Violations it left Unfixed.
type Changes struct {
	Fixed   []Path
	Removed []Path
	Unfixed []Violation
}

// Union is both fixes' changes, these first.
func (c Changes) Union(other Changes) Changes {
	return Changes{
		Fixed:   append(c.Fixed, other.Fixed...),
		Removed: append(c.Removed, other.Removed...),
		Unfixed: append(c.Unfixed, other.Unfixed...),
	}
}

// fixOf is the change a fix of path made, or its error.
func fixOf(path Path, err error) (Changes, error) {
	if err == nil {
		return Changes{Fixed: []Path{path.Printable()}}, nil
	} else {
		return Changes{}, err
	}
}

func violationsAt(path Path, messages ...string) []Violation {
	var violations []Violation
	for _, message := range messages {
		violations = append(violations, Violation{Path: path.Printable(), Message: message})
	}
	return violations
}
