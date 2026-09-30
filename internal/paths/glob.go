// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package paths

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var wildcards = strings.NewReplacer(`/\*\*/`, "/(?:.*/)?", `\*\*/`, "(?:.*/)?", `/\*\*`, "(?:/.*)?", `\*\*`, ".*", `\*`, "[^/]*", `\?`, "[^/]")

func Matcher(globs []Glob) func(Path) bool {
	var patterns []*regexp.Regexp
	for _, glob := range globs {
		patterns = append(patterns, regexp.MustCompile("^"+wildcards.Replace(regexp.QuoteMeta(cleanPath(string(glob))))+"$"))
	}
	return func(path Path) bool {
		return slices.ContainsFunc(patterns, func(pattern *regexp.Regexp) bool { return pattern.MatchString(cleanPath(string(path))) })
	}
}

func cleanPath(path string) string {
	return strings.TrimPrefix(filepath.ToSlash(filepath.Clean(path)), "./")
}
