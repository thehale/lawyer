// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import "strings"

type BlockComment struct {
	Open, Close         string
	Before, Each, After string
	Markers             []string
}

func (b BlockComment) Read(lines []string) []Line {
	return commentLines(lines[:leadingCount(lines, b.Markers, []BlockComment{b})], b.Markers, []BlockComment{b})
}

func (b BlockComment) Text(line string) string {
	return textOf(strings.TrimSpace(line), b.Markers, []BlockComment{b})
}

func (b BlockComment) Render(lines []string) []string {
	var comment []string
	if b.Before != "" {
		comment = append(comment, b.Before)
	}
	for _, line := range lines {
		comment = append(comment, b.Each+line)
	}
	if b.After != "" {
		comment = append(comment, b.After)
	}
	return comment
}

func (b BlockComment) content(line string) string {
	for _, opening := range []string{b.Before, b.Open} {
		if opening != "" {
			line = strings.TrimPrefix(line, opening)
		}
	}
	for _, closing := range []string{strings.TrimSpace(b.After), b.Close} {
		if closing != "" {
			line = strings.TrimSuffix(line, closing)
		}
	}
	if prefix := strings.TrimSpace(b.Each); prefix != "" {
		line = strings.TrimPrefix(strings.TrimSpace(line), prefix)
	}
	return strings.TrimSpace(line)
}
