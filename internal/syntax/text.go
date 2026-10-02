// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import (
	"strings"
	"unicode"
)

func HasWords(text string) bool {
	return strings.ContainsFunc(text, func(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) })
}

// ByteOrderMark is the mark some files start with, before their first line.
const ByteOrderMark = "\ufeff"
