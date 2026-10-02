// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import "regexp"

// A Prolog is a kind of line that must stay above the header. It says how many
// of the lines at the top of the file, past those already kept, it keeps.
type Prolog func(top Top) int

// Top is the top of a file: its lines, and how many of them are kept above
// the header so far.
type Top struct {
	lines []string
	kept  int
}

func (t Top) Rest() []string {
	return t.lines[min(t.kept, len(t.lines)):]
}

func leadingLine(pattern string) Prolog {
	matches := regexp.MustCompile(pattern)
	return func(top Top) int {
		if rest := top.Rest(); len(rest) > 0 && matches.MatchString(rest[0]) {
			return 1
		} else {
			return 0
		}
	}
}

var Shebang = leadingLine(`^#!($|[^\[])`)
