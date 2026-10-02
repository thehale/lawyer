// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

import (
	"regexp"
	"slices"
	"strings"

	"github.com/thehale/lawyer/internal/detect"
)

// A Prolog is a kind of line that must stay above the header. It says how many
// of the lines at the top of the file, past those already kept, it keeps.
type Prolog func(top Top) int

// Top is the top of a file: its lines, and how many of them are kept above
// the header so far.
type Top struct {
	lines []string
	kept  int
}

func (t Top) Rest() []string {
	return t.lines[min(t.kept, len(t.lines)):]
}

func leadingLine(pattern string) Prolog {
	matches := regexp.MustCompile(pattern)
	return func(top Top) int {
		if rest := top.Rest(); len(rest) > 0 && matches.MatchString(rest[0]) {
			return 1
		} else {
			return 0
		}
	}
}

var (
	Shebang     = leadingLine(`^#!($|[^\[])`)
	Declaration = leadingLine(`^\s*<\?xml`)
	Doctype     = leadingLine(`(?i)^\s*<!doctype`)
	Charset     = leadingLine(`^@charset\s`)
	OpeningTag  = leadingLine(`^<\?php`)
	Rack        = leadingLine(`^#\\`)
	Modeline    = leadingLine(`-\*-.*-\*-`)
)

var encoding = regexp.MustCompile(`^[ \t\f]*#.*?coding[:=][ \t]*[-\w.]+`)

func Encoding(top Top) int {
	encodingAt := slices.IndexFunc(top.lines[:min(2, len(top.lines))], encoding.MatchString)
	return max(0, encodingAt+1-top.kept)
}

func FrontMatter(top Top) int {
	rest := top.Rest()
	fence := ""
	if len(rest) > 0 && slices.Contains([]string{"---", "+++"}, strings.TrimSpace(rest[0])) {
		fence = strings.TrimSpace(rest[0])
	}
	closing := slices.IndexFunc(rest[min(1, len(rest)):], func(line string) bool {
		return strings.TrimSpace(line) == fence
	})

	if fence != "" && closing >= 0 {
		return closing + 2
	} else {
		return 0
	}
}

var magicComment = regexp.MustCompile(`^#\s*(?:-\*-\s*)?(?:frozen_string_literal|warn_indent|shareable_constant_value|typed):`)

func MagicComments(top Top) int {
	rest := top.Rest()
	comments := leadingMatches(rest, magicComment)

	if comments > 0 && comments < len(rest) && strings.TrimSpace(rest[comments]) == "" {
		return comments + 1
	} else {
		return comments
	}
}

var directive = regexp.MustCompile(`(?i)^#\s*(?:syntax|escape|check)\s*=`)

func Directives(top Top) int {
	return leadingMatches(top.Rest(), directive)
}

func Environment(top Top) int {
	block := docblock(top.Rest())

	if strings.Contains(strings.Join(block, "\n"), "@jest-environment") {
		return len(block)
	} else {
		return 0
	}
}

func docblock(lines []string) []string {
	closing := slices.IndexFunc(lines, func(line string) bool {
		return strings.Contains(line, "*/")
	})

	if len(lines) > 0 && strings.HasPrefix(lines[0], "/**") && closing >= 0 {
		return lines[:closing+1]
	} else {
		return nil
	}
}

var remark = regexp.MustCompile(`^\s*(?:#|//)`)

func Policy(top Top) int {
	rest := top.Rest()
	description := leadingMatches(rest, remark)
	blankAfter := description > 0 && description < len(rest) && strings.TrimSpace(rest[description]) == ""

	if blankAfter && !hasLicenseClaim(rest[:description]) {
		return description + 1
	} else {
		return 0
	}
}

func hasLicenseClaim(lines []string) bool {
	var texts []string
	for _, line := range lines {
		texts = append(texts, Mixed.Text(line))
	}
	return slices.ContainsFunc(detect.ClaimsIn(texts), func(claim detect.Claim) bool {
		switch claim.(type) {
		case detect.Identifier, detect.Notice:
			return true
		default:
			return false
		}
	})
}

func leadingMatches(lines []string, pattern *regexp.Regexp) int {
	count := 0
	for count < len(lines) && pattern.MatchString(lines[count]) {
		count++
	}
	return count
}
