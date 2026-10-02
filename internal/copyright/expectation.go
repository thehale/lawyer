// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package copyright

import "github.com/thehale/lawyer/internal/years"

// An Expectation is the copyright line a file must carry: whose, and for which
// years.
type Expectation struct {
	Owner string
	Years years.Requirement
}

// Violations are what's wrong with the years the statement gives.
func (e Expectation) Violations(existing Statement) []string {
	return e.YearViolations(existing.Years())
}

// YearViolations are what's wrong with years as written, or with their being
// unreadable.
func (e Expectation) YearViolations(current years.Range, readErr error) []string {
	switch {
	case readErr != nil:
		return []string{readErr.Error()}
	case e.Years.Violation(current) != nil:
		return []string{e.Years.Violation(current).Error()}
	default:
		return nil
	}
}
