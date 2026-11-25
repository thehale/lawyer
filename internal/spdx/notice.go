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
	ID       ID
	matcher  *regexp.Regexp
	required string
}

var StandardNotices = sync.OnceValue(func() []StandardNotice {
	files, _ := fs.Glob(texts, "texts/*.notice.txt")
	var notices []StandardNotice
	for _, file := range files {
		text, _ := texts.ReadFile(file)
		template := templateOf(string(text)).WithOptionalCopyright()
		id := ID(strings.TrimSuffix(path.Base(file), ".notice.txt"))
		notices = append(notices, StandardNotice{ID: id, matcher: template.Searcher(), required: fold(template.LongestWords())})
	}
	return notices
})

// A Text is lines run together as their words, as a notice is matched
// against them.
type Text struct {
	words, folded string
	starts        []int
}

// TextOf runs lines together as their words, keeping where each line starts.
func TextOf(lines []string) Text {
	var words strings.Builder
	starts := make([]int, len(lines))
	for index, line := range lines {
		starts[index] = words.Len()
		if normal := wordsOf(line); normal != "" {
			words.WriteString(normal)
			words.WriteString(" ")
		}
	}
	return Text{words: words.String(), folded: fold(words.String()), starts: starts}
}

// LinesIn are where the notice appears in text, from its first line up to its
// last, however its words are wrapped.
func (n StandardNotice) LinesIn(text Text) (first, last int, found bool) {
	bounds := n.boundsIn(text)

	if bounds != nil {
		start, end := words(text.words, bounds)
		return lineOf(text.starts, start), lineOf(text.starts, end) + 1, true
	} else {
		return 0, 0, false
	}
}

func (n StandardNotice) boundsIn(text Text) []int {
	if strings.Contains(text.folded, n.required) {
		return n.matcher.FindStringIndex(text.words)
	} else {
		return nil
	}
}

func words(text string, bounds []int) (int, int) {
	excerpt := text[bounds[0]:bounds[1]]
	first := bounds[0] + len(excerpt) - len(strings.TrimLeft(excerpt, " "))
	last := bounds[0] + len(strings.TrimRight(excerpt, " ")) - 1
	return first, last
}

func lineOf(starts []int, offset int) int {
	line := 0
	for line+1 < len(starts) && starts[line+1] <= offset {
		line++
	}
	return line
}
