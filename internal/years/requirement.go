// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package years

import "fmt"

// A Requirement is what written years must be, and what a fix writes instead.
type Requirement interface {
	Violation(written Range) error
	Replacement(written Range, readErr error) Range
}

type past struct {
	fallback Range
}

// InThePast accepts any years not in the future. A fix keeps readable past
// years and writes the fallback otherwise.
func InThePast(fallback Range) Requirement {
	return past{fallback: fallback}
}

func (p past) Violation(current Range) error {
	if current.IsPast() {
		return nil
	} else {
		return fmt.Errorf("years %s are in the future", current)
	}
}

func (p past) Replacement(current Range, readErr error) Range {
	if readErr == nil && current.IsPast() {
		return current
	} else {
		return p.fallback
	}
}
