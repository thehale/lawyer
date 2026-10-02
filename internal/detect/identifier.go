// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package detect

import (
	"regexp"

	"github.com/thehale/lawyer/internal/spdx"
)

// An Identifier is an SPDX-License-Identifier line.
type Identifier struct {
	span
	License spdx.Expression
}

var identifierLine = regexp.MustCompile(`^SPDX-License-Identifier:\s*(.*)$`)

func init() {
	Register(func(lines []string) []Claim {
		var claims []Claim
		for index, line := range lines {
			if match := identifierLine.FindStringSubmatch(line); match != nil {
				claims = append(claims, Identifier{span: span{First: index, Last: index + 1}, License: spdx.Expression(match[1])})
			}
		}
		return claims
	})
}
