// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package paths

import (
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// A Path is the path of a file or directory.
type Path string

func (p Path) Printable() Path {
	if strings.ContainsFunc(string(p), func(r rune) bool { return !unicode.IsPrint(r) }) {
		return Path(strconv.Quote(string(p)))
	} else {
		return p
	}
}

func (p Path) Base() string {
	return filepath.Base(string(p))
}
