// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package paths

import (
	"os"
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

func (p Path) Child(name string) Path {
	return Path(filepath.Join(string(p), name))
}

func (p Path) IsDirectory() bool {
	info, err := os.Stat(string(p))
	return err == nil && info.IsDir()
}

func (p Path) Base() string {
	return filepath.Base(string(p))
}
