// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package detect

import "github.com/thehale/lawyer/internal/copyright"

// A Copyright is a copyright line.
type Copyright struct {
	span
	Line      string
	Statement copyright.Statement
}

func init() {
	Register(func(lines []string) []Claim {
		var claims []Claim
		for index, line := range lines {
			if statement, ok := copyright.StatementIn(line); ok {
				claims = append(claims, Copyright{span: span{First: index, Last: index + 1}, Line: line, Statement: statement})
			}
		}
		return claims
	})
}
