// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package header

import (
	"slices"
	"strings"

	"github.com/thehale/lawyer/internal/syntax"
)

type document struct {
	bom, ending string
	lines       []string
}

func documentOf(content string) document {
	body, hasMark := strings.CutPrefix(content, syntax.ByteOrderMark)
	lines := strings.Split(body, "\n")
	ending := ""
	if strings.HasSuffix(lines[0], "\r") {
		ending = "\r"
	}

	if hasMark {
		return document{bom: syntax.ByteOrderMark, ending: ending, lines: lines}
	} else {
		return document{ending: ending, lines: lines}
	}
}

func (d document) String() string {
	return d.bom + strings.Join(d.lines, "\n")
}

func (d document) with(index int, lines ...string) document {
	withEndings := make([]string, len(lines))
	for i, line := range lines {
		withEndings[i] = strings.TrimSuffix(line, "\r") + d.ending
	}
	d.lines = slices.Insert(slices.Clone(d.lines), index, withEndings...)
	return d
}

func (d document) withReplaced(start, end int, lines []string) document {
	d.lines = slices.Concat(d.lines[:start], lines, d.lines[end:])
	return d
}

func (d document) withInserted(index int, block []string) document {
	if index < len(d.lines) && strings.TrimSpace(d.lines[index]) != "" {
		block = append(block, "")
	}
	return d.with(index, block...)
}
