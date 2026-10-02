// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package header

import (
	"slices"
	"strings"

	"github.com/thehale/lawyer/internal/syntax"
)

// A region is a file's leading comment, where its header lives.
type region struct {
	start int
	lines []regionLine
}

type regionLine struct {
	raw          string
	syntax, ours bool
}

func (r region) end() int {
	return r.start + len(r.lines)
}

func (r region) withoutTrailingBlanks() region {
	end := len(r.lines)
	for end > 0 && r.lines[end-1].isBlank() {
		end--
	}
	return region{start: r.start, lines: r.lines[:end]}
}

func (r region) isExactly(block []string) bool {
	stray := slices.ContainsFunc(r.lines[min(len(block), len(r.lines)):], func(line regionLine) bool { return line.ours })
	return !stray && slices.Equal(r.raws(len(block)), block)
}

func (r region) raws(count int) []string {
	var raws []string
	for _, line := range r.lines[:min(count, len(r.lines))] {
		raws = append(raws, strings.TrimSuffix(line.raw, "\r"))
	}
	return raws
}

func (r region) remainingLines() []string {
	removals := r.removals()
	var kept []string
	for index, line := range r.lines {
		if !removals[index] {
			kept = append(kept, line.raw)
		}
	}

	if slices.ContainsFunc(kept, syntax.HasWords) {
		return kept
	} else {
		return nil
	}
}

func (r region) removals() []bool {
	removals := make([]bool, len(r.lines))
	for index, line := range r.lines {
		removals[index] = line.ours
	}
	for index := 1; index < len(r.lines); index++ {
		removals[index] = removals[index] || removals[index-1] && r.lines[index].isBare()
	}
	for index := len(r.lines) - 2; index >= 0; index-- {
		removals[index] = removals[index] || removals[index+1] && r.lines[index].isBare()
	}
	return removals
}

func (l regionLine) isBare() bool {
	return !l.syntax && !syntax.HasWords(l.raw)
}

func (l regionLine) isBlank() bool {
	return strings.TrimSpace(l.raw) == ""
}
