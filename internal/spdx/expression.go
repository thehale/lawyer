// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import "regexp"

var wellFormed = regexp.MustCompile(`^[(\s]*[\w.+-]+(?:[()\s]+(?:AND|OR|WITH)[(\s]+[\w.+-]+)*[)\s]*$`)

// An Expression is an SPDX license expression, such as "MPL-2.0" or
// "MIT OR Apache-2.0".
type Expression string

// IsWellFormed reports whether the expression follows SPDX's grammar, whether or
// not lawyer carries the licenses it names.
func (e Expression) IsWellFormed() bool {
	return wellFormed.MatchString(string(e))
}
