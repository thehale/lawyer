// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import "strings"

type LineComment struct {
	Prefix  string
	Markers []string
	Blocks  []BlockComment
}

func (l LineComment) Read(lines []string) []Line {
	return commentLines(lines[:leadingCount(lines, l.Markers, l.Blocks)], l.Markers, l.Blocks)
}

func (l LineComment) Text(line string) string {
	return textOf(strings.TrimSpace(line), l.Markers, l.Blocks)
}

func (l LineComment) Render(lines []string) []string {
	var comment []string
	for _, line := range lines {
		comment = append(comment, l.Prefix+line)
	}
	return comment
}
