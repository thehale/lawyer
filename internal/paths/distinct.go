// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package paths

import (
	"slices"
	"strings"
)

// Distinct are the paths sorted, each named once however it was spelled.
func Distinct(paths []Path) []Path {
	inOrder := slices.SortedFunc(slices.Values(paths), func(a, b Path) int {
		return strings.Compare(cleanPath(string(a)), cleanPath(string(b)))
	})
	return slices.CompactFunc(inOrder, func(a, b Path) bool { return cleanPath(string(a)) == cleanPath(string(b)) })
}
