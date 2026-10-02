// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

// A Violation is one way a file falls short of a Declaration.
type Violation struct {
	Path    Path
	Message string
}

func violationsAt(path Path, messages ...string) []Violation {
	var violations []Violation
	for _, message := range messages {
		violations = append(violations, Violation{Path: path.Printable(), Message: message})
	}
	return violations
}
