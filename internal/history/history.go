// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package history

import (
	"slices"
	"strings"

	"github.com/thehale/lawyer/internal/git"
	"github.com/thehale/lawyer/internal/paths"
	"github.com/thehale/lawyer/internal/years"
)

type History struct {
	root paths.Path
	// Gap is why the history can't date every edit.
	Gap string
}

func Of(root paths.Path) History {
	history := History{root: root}

	switch {
	case !git.IsInCheckout(string(root)):
		history.Gap = "this is not a git checkout"
	case history.output("rev-parse", "--is-shallow-repository") == "true":
		history.Gap = "this clone is shallow"
	}
	return history
}

func (h History) IsEveryEditDated() bool {
	return h.Gap == ""
}

func (h History) EditYears(path paths.Path) years.Range {
	if h.IsEveryEditDated() {
		return rangeOf(append(h.committedYears(path), h.uncommittedYears(path)...))
	} else {
		return years.This()
	}
}

func (h History) Years() years.Range {
	if h.IsEveryEditDated() {
		return rangeOf(append(yearsOf(h.output("log", "--format=%ad", "--date=format:%Y")), h.dirtyYears()...))
	} else {
		return years.This()
	}
}

func rangeOf(edits []int) years.Range {
	if len(edits) == 0 {
		return years.This()
	} else {
		return years.Range{First: slices.Min(edits), Last: slices.Max(edits)}
	}
}

func (h History) dirtyYears() []int {
	if h.output("status", "--porcelain") != "" {
		return []int{years.This().Last}
	} else {
		return nil
	}
}

func (h History) output(args ...string) string {
	output, _ := git.Run(string(h.root), args...)
	return strings.TrimSpace(output)
}

func (h History) committedYears(path paths.Path) []int {
	log := h.output("log", "--follow", "--format=%x00%ad", "--date=format:%Y", "--patch", "--cc", "--unified=0", "--text", "--no-color", "--no-ext-diff", "--", string(path))
	var edits []int
	for _, commit := range strings.Split(log, marker)[1:] {
		year, diff, _ := strings.Cut(commit, "\n")
		if isSubstantive(diff) {
			edits = append(edits, yearsOf(year)...)
		}
	}
	return edits
}

func (h History) uncommittedYears(path paths.Path) []int {
	diff := h.output("diff", "HEAD", "--unified=0", "--text", "--no-color", "--no-ext-diff", "--", string(path))

	if isSubstantive(diff) || h.output("ls-files", "--", string(path)) == "" {
		return []int{years.This().Last}
	} else {
		return nil
	}
}
