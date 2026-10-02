// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package detect

import (
	"cmp"
	"slices"
)

// A Claim is what a run of comment lines says about a file's copyright or
// license.
type Claim interface {
	HasLine(index int) bool
	// IsComplete reports whether the claim has every part its kind needs, such
	// as a copyright's years.
	IsComplete() bool
	start() int
}

type span struct {
	First, Last int
}

func (s span) HasLine(index int) bool {
	return s.First <= index && index < s.Last
}

func (s span) start() int {
	return s.First
}

type Pattern func(lines []string) []Claim

var patterns []Pattern

func Register(pattern Pattern) {
	patterns = append(patterns, pattern)
}

// Claims are claims found together in one run of lines.
type Claims []Claim

// ClaimsIn are the claims every pattern finds in lines, in the order they
// start.
func ClaimsIn(lines []string) Claims {
	var claims Claims
	for _, pattern := range patterns {
		claims = append(claims, pattern(lines)...)
	}
	slices.SortStableFunc(claims, func(a, b Claim) int { return cmp.Compare(a.start(), b.start()) })
	return claims
}

func (c Claims) HasLine(index int) bool {
	return slices.ContainsFunc(c, func(claim Claim) bool { return claim.HasLine(index) })
}
