// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"regexp"
	"slices"
	"strings"
)

// A License is one license's SPDX text, and the template its copies are read
// against.
type License struct {
	ID       ID
	Text     string
	template template
}

// HasCopyrightSlot reports whether the license's text has a place for a
// copyright line, such as MIT's "Copyright (c) <year> <copyright holders>".
func (l License) HasCopyrightSlot() bool {
	return l.template.Slot() != nil
}

// IsWrittenIn reports whether content is the license's text, however it is
// wrapped or formatted, with any copyright line as written.
func (l License) IsWrittenIn(content string) bool {
	return l.template.IsMatch(wordsOf(content))
}

// OwnerYears are the years written on owner's copyright line in content, such
// as "2020-2024".
func (l License) OwnerYears(content, owner string) (string, bool) {
	line := l.template.CopyrightGroups(ownerLine(owner), wordsOf(content))

	if line != nil {
		return strings.TrimSuffix(line["first"]+"-"+line["last"], "-"), true
	} else {
		return "", false
	}
}

// HasCopyrightLine reports whether content has a copyright line, whoever's,
// where the license takes one.
func (l License) HasCopyrightLine(content string) bool {
	return l.template.CopyrightGroups("(?:.*?)", wordsOf(content)) != nil
}

func ownerLine(owner string) string {
	return `copyright c (?P<first>\d{4})(?: (?P<last>\d{4}))? ` + regexp.QuoteMeta(wordsOf(owner))
}

// TextWithCopyright is the license's text with line in its copyright slot, or
// the text alone for a license without one.
func (l License) TextWithCopyright(line string) string {
	lines := strings.Split(l.Text, "\n")
	slot := l.template.Slot()
	first, last := placeholder(lines, wordsOf(original(slot)))

	switch {
	case first >= 0:
		return strings.Join(slices.Concat(lines[:first], []string{line}, lines[last:]), "\n")
	case slot != nil:
		return line + "\n\n" + l.Text
	default:
		return l.Text
	}
}

func placeholder(lines []string, words string) (int, int) {
	for first := range lines {
		for last := first + 1; last <= min(first+6, len(lines)); last++ {
			if words != "" && strings.TrimSpace(lines[first]) != "" && wordsOf(strings.Join(lines[first:last], " ")) == words {
				return first, last
			}
		}
	}
	return -1, -1
}

func original(slot *variable) string {
	if slot == nil {
		return ""
	} else {
		return slot.original
	}
}
