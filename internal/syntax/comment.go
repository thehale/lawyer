// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import (
	"slices"
	"strings"
)

type Comment interface {
	Read(lines []string) []Line
	Text(line string) string
	Render(lines []string) []string
}

type Line struct {
	Text   string
	Syntax bool
}

func leadingCount(lines []string, markers []string, blocks []BlockComment) int {
	var open *BlockComment
	count := 0
	for count < len(lines) && (open != nil || isCommentStart(lines[count], markers, blocks)) {
		open = openBlock(lines[count], open, blocks)
		count++
	}
	return count
}

func isCommentStart(line string, markers []string, blocks []BlockComment) bool {
	text := strings.ToLower(strings.TrimSpace(line))
	hasMarker := slices.ContainsFunc(markers, func(marker string) bool { return strings.HasPrefix(text, marker) })
	return text == "" || hasMarker || opening(text, blocks) != nil
}

func opening(text string, blocks []BlockComment) *BlockComment {
	index := slices.IndexFunc(blocks, func(block BlockComment) bool { return strings.HasPrefix(text, block.Open) })

	if index >= 0 {
		return &blocks[index]
	} else {
		return nil
	}
}

func openBlock(line string, open *BlockComment, blocks []BlockComment) *BlockComment {
	text := strings.TrimSpace(line)
	block := opening(text, blocks)

	switch {
	case open != nil && strings.Contains(text, open.Close):
		return nil
	case open != nil:
		return open
	case block != nil && !strings.Contains(text[len(block.Open):], block.Close):
		return block
	default:
		return nil
	}
}

func commentLines(lines []string, markers []string, blocks []BlockComment) []Line {
	comment := make([]Line, len(lines))
	for index, line := range lines {
		text := textOf(strings.TrimSpace(line), markers, blocks)
		comment[index] = Line{Text: text, Syntax: text == "" && hasDelimiter(line, blocks)}
	}
	return comment
}

func textOf(line string, markers []string, blocks []BlockComment) string {
	for _, marker := range markers {
		if len(line) >= len(marker) && strings.EqualFold(line[:len(marker)], marker) {
			return strings.TrimSpace(line[len(marker):])
		}
	}
	for _, block := range blocks {
		line = block.content(line)
	}
	return line
}

func hasDelimiter(line string, blocks []BlockComment) bool {
	return slices.ContainsFunc(blocks, func(block BlockComment) bool {
		return strings.Contains(line, block.Open) || strings.Contains(line, block.Close)
	})
}
