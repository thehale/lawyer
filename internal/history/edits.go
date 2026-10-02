// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package history

import (
	"slices"
	"strconv"
	"strings"

	"github.com/thehale/lawyer/internal/detect"
	"github.com/thehale/lawyer/internal/syntax"
	"github.com/thehale/lawyer/internal/years"
)

const marker = "\x00"

func isSubstantive(diff string) bool {
	var contents []string
	for _, line := range changes(diff) {
		contents = append(contents, syntax.CommentText(line))
	}
	claims := slices.DeleteFunc(detect.ClaimsIn(contents), func(claim detect.Claim) bool { return !claim.IsComplete() })
	edit := false
	for index, content := range contents {
		edit = edit || syntax.HasWords(content) && !claims.HasLine(index)
	}
	return edit
}

func changes(diff string) []string {
	var lines []string
	parents := 0
	for _, line := range strings.Split(diff, "\n") {
		switch {
		case strings.HasPrefix(line, "diff --"):
			parents = 0
		case strings.HasPrefix(line, "@@"):
			parents = len(line) - len(strings.TrimLeft(line, "@")) - 1
		case parents > 0 && len(line) >= parents && strings.ContainsAny(line[:parents], "+-"):
			lines = append(lines, line[parents:])
		}
	}
	return lines
}

func yearsOf(text string) []int {
	var dates []int
	for _, field := range strings.Fields(text) {
		if year, err := strconv.Atoi(field); err == nil && year <= years.This().Last {
			dates = append(dates, year)
		}
	}
	return dates
}
