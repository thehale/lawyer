// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"cmp"
	"slices"
	"sync"
)

// A known license is one lawyer carries, with its text's words to compare.
type known struct {
	license License
	words   string
}

// catalog is every known license, base licenses before their variants.
var catalog = sync.OnceValue(func() []known {
	var licenses []known
	for _, id := range IDs() {
		license, _ := ByID(id)
		licenses = append(licenses, known{license: license, words: wordsOf(license.Text)})
	}
	slices.SortStableFunc(licenses, func(a, b known) int { return cmp.Compare(len(a.license.ID), len(b.license.ID)) })
	return licenses
})

// Identify is the license text is: one whose text it is word for word, or
// else the one closest in length whose template it matches. Where several
// licenses share a text, as MPL-2.0 and MPL-2.0-no-copyleft-exception do, it
// is the base license rather than a variant.
func Identify(text string) (License, bool) {
	words := wordsOf(text)
	licenses := catalog()
	exact := slices.IndexFunc(licenses, func(candidate known) bool { return candidate.words == words })

	if exact >= 0 {
		return licenses[exact].license, true
	} else {
		return closestMatch(nearInLength(licenses, len(words)), words)
	}
}

func closestMatch(near []known, words string) (License, bool) {
	index := slices.IndexFunc(near, func(candidate known) bool { return candidate.license.template.IsMatch(words) })

	if index >= 0 {
		return near[index].license, true
	} else {
		return License{}, false
	}
}

// nearInLength are the licenses within a quarter of length words of it,
// closest first, so a tolerant match is tried against few full templates.
func nearInLength(licenses []known, length int) []known {
	distance := func(candidate known) int { return max(len(candidate.words)-length, length-len(candidate.words)) }
	near := slices.DeleteFunc(slices.Clone(licenses), func(candidate known) bool { return distance(candidate) > length/4 })
	slices.SortStableFunc(near, func(a, b known) int { return cmp.Compare(distance(a), distance(b)) })
	return near
}
