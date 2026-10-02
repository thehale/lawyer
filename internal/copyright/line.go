// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package copyright

import "github.com/thehale/lawyer/internal/years"

// A Line is a copyright line in lawyer's form, "Copyright (c) <years> <owner>".
type Line struct {
	Years years.Range
	Owner string
}

func (l Line) String() string {
	return "Copyright (c) " + l.Years.String() + " " + l.Owner
}
