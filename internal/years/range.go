// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package years

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Range struct {
	First, Last int
}

func Only(year int) Range {
	return Range{First: year, Last: year}
}

func This() Range {
	return Only(time.Now().Year())
}

func RangeOf(text string) (Range, error) {
	first, last, isRange := strings.Cut(text, "-")
	if !isRange {
		last = first
	}
	from, fromErr := strconv.Atoi(first)
	to, toErr := strconv.Atoi(last)

	if fromErr == nil && toErr == nil && len(first) == 4 && len(last) == 4 && from <= to {
		return Range{First: from, Last: to}, nil
	} else {
		return Range{}, fmt.Errorf("unreadable years %q", text)
	}
}

func (r Range) String() string {
	if r.First == r.Last {
		return strconv.Itoa(r.First)
	} else {
		return fmt.Sprintf("%d-%d", r.First, r.Last)
	}
}

func (r Range) IsPast() bool {
	return r.Last <= time.Now().Year()
}
