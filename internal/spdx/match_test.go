// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"regexp"
	"strings"
	"testing"
)

var punctuation = regexp.MustCompile(`[^\pL\pN]+`)

func FuzzWordsOfSplitsAsTheRegexpDoes(f *testing.F) {
	for _, seed := range []string{"", " ", "Hello, World!", "© 2024 Joseph Hale", "see https://x and éhttpsé", "\xff\xfeab  c\t", "x²³ ½ Ⅻ ٣", "İSTANBUL ǅ ﬁ", "a_b-c.d"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, text string) {
		lowered := strings.ReplaceAll(strings.ToLower(text), "©", " c ")
		want := secure.ReplaceAllString(strings.TrimSpace(punctuation.ReplaceAllString(lowered, " ")), "http")
		if got := wordsOf(text); got != want {
			t.Errorf("wordsOf(%q) = %q, want %q", text, got, want)
		}
	})
}
