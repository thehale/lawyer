// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"regexp"
	"slices"
	"strings"
)

// A template is an SPDX license template: literal text, optional groups, and
// variables standing for text that may vary, one of them perhaps the copyright
// line.
type template struct {
	parts []part
}

type part interface {
	source(slot *variable, copyright string) string
}

type literal struct {
	words string
}

type optional struct {
	parts []part
}

type variable struct {
	name, original, match string
}

var (
	markup       = regexp.MustCompile(`(?s)<<(beginOptional|endOptional|var;.*?")>>`)
	attribute    = regexp.MustCompile(`(?:^var|");(\w+)="`)
	placeholders = regexp.MustCompile(`<[^>]+>|\[[^\]]+\]`)
)

func templateOf(text string) template {
	levels := [][]part{nil}
	pieces := markup.Split(text, -1)
	for index, tag := range markup.FindAllStringSubmatch(text, -1) {
		levels[len(levels)-1] = append(levels[len(levels)-1], literal{words: wordsOf(pieces[index])})
		levels = levelsAfter(levels, tag[1])
	}
	return template{parts: append(levels[0], literal{words: wordsOf(pieces[len(pieces)-1])})}
}

func levelsAfter(levels [][]part, tag string) [][]part {
	top := len(levels) - 1

	switch tag {
	case "beginOptional":
		return append(levels, nil)
	case "endOptional":
		levels[top-1] = append(levels[top-1], optional{parts: levels[top]})
		return levels[:top]
	default:
		levels[top] = append(levels[top], variableIn(tag))
		return levels
	}
}

func variableIn(tag string) *variable {
	attributes := map[string]string{}
	starts := attribute.FindAllStringSubmatchIndex(tag, -1)
	for index, start := range starts {
		ends := len(tag) - 1
		if index+1 < len(starts) {
			ends = starts[index+1][0]
		}
		attributes[tag[start[2]:start[3]]] = tag[start[1]:ends]
	}
	return &variable{name: attributes["name"], original: attributes["original"], match: attributes["match"]}
}

// Slot is the variable a copyright line takes the place of: a top-level
// copyright variable whose original text is a placeholder, such as
// "<year> <owner>". A concrete notice, such as the FSF's, is no slot.
func (t template) Slot() *variable {
	index := slices.IndexFunc(t.parts, func(each part) bool {
		candidate, ok := each.(*variable)
		return ok && candidate.name == "copyright" && placeholders.MatchString(candidate.original)
	})

	if index >= 0 {
		return t.parts[index].(*variable)
	} else {
		return nil
	}
}

// IsMatch reports whether text is the template's own text, with any copyright
// variable as written.
func (t template) IsMatch(text string) bool {
	return wholeMatcher(t.source(nil, "")).MatchString(text)
}

// CopyrightGroups are the named groups of the copyright line copyright matches
// in text: in the slot if there is one, and otherwise on top of the template's
// text.
func (t template) CopyrightGroups(copyright, text string) map[string]string {
	slot := t.Slot()
	slotMatcher := wholeMatcher(t.source(slot, copyright))

	if slot != nil && slotMatcher.MatchString(text) {
		return groups(slotMatcher, text)
	} else {
		return groups(wholeMatcher(copyright+" ?"+t.source(nil, "")), text)
	}
}

func groups(matcher *regexp.Regexp, text string) map[string]string {
	match := matcher.FindStringSubmatch(text)

	if match != nil {
		named := make(map[string]string, len(match))
		for group, name := range matcher.SubexpNames() {
			named[name] = match[group]
		}
		return named
	} else {
		return nil
	}
}

// Searcher finds the template's text anywhere within a longer text.
func (t template) Searcher() *regexp.Regexp {
	return regexp.MustCompile("(?i)(?:^| )" + t.source(nil, "") + "(?: |$)")
}

// WithOptionalCopyright makes a top-level copyright variable optional, along
// with the "Copyright" or "Copyright (c)" before it.
func (t template) WithOptionalCopyright() template {
	var parts []part
	for _, each := range t.parts {
		parts = withOptional(parts, each)
	}
	return template{parts: parts}
}

func withOptional(parts []part, next part) []part {
	copyright, isVariable := next.(*variable)
	before, follows := last[literal](parts)
	key := copyrightKey(before.words)

	if isVariable && copyright.name == "copyright" && follows && key != "" {
		prefix := strings.TrimSpace(strings.TrimSuffix(before.words, key))
		return append(parts[:len(parts)-1], literal{words: prefix}, optional{parts: []part{literal{words: key}, copyright}})
	} else {
		return append(parts, next)
	}
}

func copyrightKey(words string) string {
	for _, key := range []string{"copyright c", "copyright"} {
		if words == key || strings.HasSuffix(words, " "+key) {
			return key
		}
	}
	return ""
}

func last[T part](parts []part) (T, bool) {
	var zero T
	if len(parts) == 0 {
		return zero, false
	} else {
		typed, ok := parts[len(parts)-1].(T)
		return typed, ok
	}
}

func (t template) source(slot *variable, copyright string) string {
	return sourceOf(t.parts, slot, copyright)
}

func sourceOf(parts []part, slot *variable, copyright string) string {
	var source strings.Builder
	for _, each := range parts {
		source.WriteString(each.source(slot, copyright))
	}
	return source.String()
}

func (l literal) source(*variable, string) string {
	if l.words == "" {
		return ""
	} else {
		return " ?" + regexp.QuoteMeta(l.words) + " ?"
	}
}

func (o optional) source(slot *variable, copyright string) string {
	return "(?:" + sourceOf(o.parts, slot, copyright) + ")?"
}

func (v *variable) source(slot *variable, copyright string) string {
	if v == slot {
		return " ?" + copyright + " ?"
	} else {
		return " ?(?:" + cappedPattern(tolerantPattern(v.match)) + ")? ?"
	}
}

func wholeMatcher(source string) *regexp.Regexp {
	return regexp.MustCompile("(?i)^ ? ?" + source + " ?$")
}
