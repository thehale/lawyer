// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package history

import (
	"runtime"
	"sync"

	"github.com/thehale/lawyer/internal/paths"
	"github.com/thehale/lawyer/internal/years"
)

func (h History) Date(files []paths.Path) History {
	if h.IsEveryEditDated() {
		h.dated = h.editYearsOf(files, h.worktree())
	}
	return h
}

func (h History) editYearsOf(files []paths.Path, tree worktree) map[paths.Path]years.Range {
	ranges := make([]years.Range, len(files))
	slots := make(chan struct{}, runtime.NumCPU())
	var group sync.WaitGroup
	for index, path := range files {
		slots <- struct{}{}
		group.Go(func() {
			ranges[index] = rangeOf(append(h.committedYears(path), tree.uncommittedYears(path)...))
			<-slots
		})
	}
	group.Wait()
	return mapOf(files, ranges)
}

func mapOf[K comparable, V any](keys []K, values []V) map[K]V {
	pairs := make(map[K]V, len(keys))
	for index, key := range keys {
		pairs[key] = values[index]
	}
	return pairs
}
