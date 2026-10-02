// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package years

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
)

var (
	yearList = regexp.MustCompile(`^[0-9]{4}(?:\s*[-–,]\s*[0-9]{4})*$`)
	year     = regexp.MustCompile(`[0-9]{4}`)
)

func LooseRange(text string) (Range, error) {
	if yearList.MatchString(text) {
		return span(year.FindAllString(text, -1)), nil
	} else {
		return Range{}, fmt.Errorf("unreadable years %q", text)
	}
}

func span(digits []string) Range {
	var numbers []int
	for _, year := range digits {
		number, _ := strconv.Atoi(year)
		numbers = append(numbers, number)
	}
	return Range{First: slices.Min(numbers), Last: slices.Max(numbers)}
}
