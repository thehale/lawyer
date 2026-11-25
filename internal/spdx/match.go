// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	repeat = regexp.MustCompile(`\{(\d+),(\d+)\}`)
	secure = regexp.MustCompile(`\bhttps\b`)

	classes            = regexp.MustCompile(`\[(?:\\.|[^\]])*\]|\{\d*,?\d*\}|(?:\\.|[^\[{])+|\{`)
	escapedPunctuation = regexp.MustCompile(`\\[^\w\s]`)
	punctuationClass   = regexp.MustCompile(`^\[(?:\\[^\w\s]|[^\w\s\]\\])+\]$`)
	literalPunctuation = regexp.MustCompile(`[-,'"":;/!@#%&=<>~` + "`" + `]`)
)

func wordsOf(text string) string {
	text = strings.ReplaceAll(strings.ToLower(text), "©", " c ")
	words := strings.Join(strings.FieldsFunc(text, isPunctuation), " ")
	return secure.ReplaceAllString(words, "http")
}

func isPunctuation(character rune) bool {
	return !unicode.IsLetter(character) && !unicode.IsNumber(character)
}

func tolerantPattern(match string) string {
	return classes.ReplaceAllStringFunc(match, func(piece string) string {
		if punctuationClass.MatchString(piece) {
			return "(?: ?)"
		} else if strings.HasPrefix(piece, "[") || strings.HasPrefix(piece, "{") {
			return piece
		} else {
			withMarkers := literalPunctuation.ReplaceAllString(escapedPunctuation.ReplaceAllString(piece, "\x00"), "\x00")
			return strings.ReplaceAll(withMarkers, "\x00", "(?: ?)")
		}
	})
}

func cappedPattern(match string) string {
	return repeat.ReplaceAllStringFunc(match, func(counted string) string {
		limits := repeat.FindStringSubmatch(counted)
		if most, _ := strconv.Atoi(limits[2]); most > 1000 {
			return "{" + limits[1] + ",}"
		} else {
			return counted
		}
	})
}
