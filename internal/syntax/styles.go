// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package syntax

var (
	Hashes     = LineComment{Prefix: "# ", Markers: []string{"#"}}
	Mixed      = LineComment{Prefix: "# ", Markers: []string{"#", "//"}, Blocks: []BlockComment{{Open: "/*", Close: "*/"}}}
	Slashes    = LineComment{Prefix: "// ", Markers: []string{"//"}, Blocks: []BlockComment{{Open: "/*", Close: "*/"}}}
	Stars      = BlockComment{Open: "/*", Close: "*/", Before: "/*", Each: " * ", After: " */", Markers: []string{"//"}}
	Markup     = BlockComment{Open: "<!--", Close: "-->", Before: "<!--", After: "-->"}
	Remarks    = LineComment{Prefix: "@rem ", Markers: []string{"@rem", "rem", "::"}}
	Dashes     = LineComment{Prefix: "-- ", Markers: []string{"--"}}
	Semicolons = LineComment{Prefix: ";; ", Markers: []string{";"}}
	Percents   = LineComment{Prefix: "%% ", Markers: []string{"%"}}
	Parens     = BlockComment{Open: "(*", Close: "*)", Before: "(*", After: "*)"}
	Mustaches  = BlockComment{Open: "{{!", Close: "}}", Before: "{{!--", After: "--}}"}
	Scriptlets = BlockComment{Open: "<%", Close: "%>", Before: "<%#", After: "%>"}
)
