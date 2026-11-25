// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package history

import (
	"path/filepath"
	"strings"

	"github.com/thehale/lawyer/internal/git"
	"github.com/thehale/lawyer/internal/paths"
)

type worktree struct {
	history History
	root    string
	changed map[string]bool
	tracked map[string]bool
}

func (h History) worktree() worktree {
	root, _ := filepath.Abs(string(h.root))
	changed, _ := git.Run(string(h.root), "diff", "HEAD", "--name-only", "-z", "--no-renames", "--relative")
	index, _ := git.Run(string(h.root), "ls-files", "-z")
	return worktree{history: h, root: root, changed: setOf(namesIn(changed)), tracked: setOf(namesIn(index))}
}

func (w worktree) uncommittedYears(path paths.Path) []int {
	name, isUnder := w.nameOf(path)

	if isUnder && !w.changed[name] {
		return thisYearIf(!w.tracked[name])
	} else {
		return w.history.uncommittedYears(path)
	}
}

func (w worktree) nameOf(path paths.Path) (string, bool) {
	name, err := filepath.Rel(w.root, w.locate(string(path)))
	return filepath.ToSlash(name), err == nil && filepath.IsLocal(name)
}

func (w worktree) locate(path string) string {
	if filepath.IsAbs(path) {
		return path
	} else {
		return filepath.Join(w.root, path)
	}
}

func namesIn(list string) []string {
	return strings.Split(strings.TrimSuffix(list, "\x00"), "\x00")
}

func setOf(members []string) map[string]bool {
	set := make(map[string]bool, len(members))
	for _, member := range members {
		set[member] = true
	}
	return set
}
