// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"strings"
	"unicode"
)

// fold replaces each letter in text with the least of the letters a
// case-insensitive regexp takes for it, so two texts that match that way
// compare equal.
func fold(text string) string {
	return strings.Map(least, text)
}

func least(letter rune) rune {
	smallest := letter
	for other := unicode.SimpleFold(letter); other != letter; other = unicode.SimpleFold(other) {
		smallest = min(smallest, other)
	}
	return smallest
}
