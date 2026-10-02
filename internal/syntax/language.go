// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import (
	"slices"
	"strings"
)

type Language struct {
	Name       string
	Filenames  []string
	Extensions []string
	Shebangs   []string
	Comment    Comment
	Prolog     []Prolog
}

func (l Language) HasHeader() bool {
	return l.Comment != nil
}

func (l Language) HeaderStart(lines []string) int {
	top := Top{lines: lines}
	for _, prolog := range l.Prolog {
		top.kept += prolog(top)
	}
	return top.kept
}

func (l Language) Text(line string) string {
	if l.HasHeader() {
		return l.Comment.Text(line)
	} else {
		return strings.TrimSpace(line)
	}
}

func (l Language) hasName(base string) bool {
	byName := slices.ContainsFunc(l.Filenames, func(name string) bool {
		return strings.EqualFold(name, base)
	})
	byExtension := slices.ContainsFunc(l.Extensions, func(extension string) bool {
		return strings.HasSuffix(base, "."+extension)
	})
	return byName || byExtension
}
