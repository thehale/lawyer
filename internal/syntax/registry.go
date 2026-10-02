// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import (
	"regexp"
	"slices"
	"strings"
)

var languages []Language

func Register(language Language) {
	languages = append(languages, language)
}

func KnownLanguages() []Language {
	return slices.SortedFunc(slices.Values(languages), func(a, b Language) int {
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	})
}

func LanguageOf(name, content string) (Language, bool) {
	base := strings.ToLower(name)
	interpreter := interpreterOf(firstLine(content))

	index := slices.IndexFunc(languages, func(language Language) bool {
		return language.hasName(base) || slices.Contains(language.Shebangs, interpreter)
	})

	if index >= 0 {
		return languages[index], true
	} else {
		return Language{}, false
	}
}

var interpreterLine = regexp.MustCompile(`^#!\s*(?:\S*/)?([^/\s]+)(?:\s+(?:-\S+\s+)*([^-\s]\S*))?`)

func interpreterOf(firstLine string) string {
	match := interpreterLine.FindStringSubmatch(firstLine)

	if match == nil {
		return ""
	} else if match[1] == "env" {
		return match[2]
	} else {
		return match[1]
	}
}

func firstLine(content string) string {
	line, _, _ := strings.Cut(strings.TrimPrefix(content, ByteOrderMark), "\n")
	return line
}

// CommentText is a line's text without the comment syntax of any language
// lawyer knows, for a line whose language isn't known.
func CommentText(line string) string {
	content := strings.TrimSpace(line)
	for _, language := range languages {
		if text := language.Text(line); language.HasHeader() && text != content {
			return text
		}
	}
	return content
}
