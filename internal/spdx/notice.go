// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"io/fs"
	"path"
	"regexp"
	"strings"
	"sync"
)

// A StandardNotice is the paragraph a license asks files to carry, such as the
// MPL's "This Source Code Form is subject to…".
type StandardNotice struct {
	ID      ID
	matcher *regexp.Regexp
}

var StandardNotices = sync.OnceValue(func() []StandardNotice {
	files, _ := fs.Glob(texts, "texts/*.notice.txt")
	var notices []StandardNotice
	for _, file := range files {
		text, _ := texts.ReadFile(file)
		matcher := templateOf(string(text)).WithOptionalCopyright().Searcher()
		notices = append(notices, StandardNotice{ID: ID(strings.TrimSuffix(path.Base(file), ".notice.txt")), matcher: matcher})
	}
	return notices
})

// LinesIn are where the notice's text appears in lines, from first up to
// last, however its words are wrapped.
func (n StandardNotice) LinesIn(lines []string) (first, last int, found bool) {
	text, starts := textOf(lines)
	bounds := n.matcher.FindStringIndex(text)

	if bounds != nil {
		start, end := words(text, bounds)
		return lineOf(starts, start), lineOf(starts, end) + 1, true
	} else {
		return 0, 0, false
	}
}

func words(text string, bounds []int) (int, int) {
	excerpt := text[bounds[0]:bounds[1]]
	first := bounds[0] + len(excerpt) - len(strings.TrimLeft(excerpt, " "))
	last := bounds[0] + len(strings.TrimRight(excerpt, " ")) - 1
	return first, last
}

func textOf(lines []string) (string, []int) {
	var text strings.Builder
	starts := make([]int, len(lines))
	for index, line := range lines {
		starts[index] = text.Len()
		if normal := wordsOf(line); normal != "" {
			text.WriteString(normal)
			text.WriteString(" ")
		}
	}
	return text.String(), starts
}

func lineOf(starts []int, offset int) int {
	line := 0
	for line+1 < len(starts) && starts[line+1] <= offset {
		line++
	}
	return line
}
