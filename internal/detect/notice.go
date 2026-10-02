// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package detect

import (
	"cmp"
	"slices"

	"github.com/thehale/lawyer/internal/spdx"
)

// A Notice is a license's standard notice, such as the paragraph the MPL asks
// files to carry.
type Notice struct {
	span
	License spdx.Expression
}

func (n Notice) IsComplete() bool {
	return true
}

func init() {
	Register(func(lines []string) []Claim {
		var notices []Notice
		for _, standard := range spdx.StandardNotices() {
			if first, last, found := standard.LinesIn(lines); found {
				notices = append(notices, Notice{span: span{First: first, Last: last}, License: standard.ID.Expression()})
			}
		}

		if notices == nil {
			return nil
		} else {
			return []Claim{slices.MinFunc(notices, func(a, b Notice) int { return cmp.Compare(a.First, b.First) })}
		}
	})
}
