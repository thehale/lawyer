// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package languages_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/thehale/lawyer/internal/copyright"
	"github.com/thehale/lawyer/internal/header"
	_ "github.com/thehale/lawyer/internal/languages"
	"github.com/thehale/lawyer/internal/syntax"
	"github.com/thehale/lawyer/internal/years"
)

var expectation = header.Expectation{Copyright: copyright.Expectation{Owner: "Joseph Hale", Years: years.InThePast(years.Only(2026))}, License: "MPL-2.0"}

func TestExamples(t *testing.T) {
	examples, _ := filepath.Glob(filepath.Join("testdata", "*", "*.txt"))
	for _, path := range examples {
		t.Run(filepath.ToSlash(path), func(t *testing.T) {
			before, after := example(t, path)
			language := languageOf(t, strings.TrimSuffix(path, ".txt"), before)
			existing := header.Read(language, before)
			if fixed := existing.ContentWith(existing.Canonical(expectation)); fixed != after {
				t.Errorf("got\n%s\nwant\n%s", fixed, after)
			}
			if violations := header.Read(language, after).Violations(expectation); violations != nil {
				t.Errorf("the fixed file still has %q", violations)
			}
		})
	}
}

var separator = regexp.MustCompile(`(?m)^-{80}\r?\n`)

func example(t *testing.T, path string) (string, string) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	bounds := separator.FindStringIndex(string(content))

	if bounds == nil {
		return string(content), string(content)
	} else {
		return string(content[:bounds[0]]), string(content[bounds[1]:])
	}
}

func TestShebangs(t *testing.T) {
	cases := map[string]string{
		"#!/usr/bin/env bash":        "Bash",
		"#!/bin/sh":                  "Bash",
		"#!/usr/bin/env -S bash -eu": "Bash",
		"#!/usr/bin/python3":         "Python",
		"#!/usr/bin/env ruby":        "Ruby",
	}
	for line, name := range cases {
		if language := languageOf(t, "bin/tool", line+"\n"); language.Name != name {
			t.Errorf("%s detected as %q, expected %q", line, language.Name, name)
		}
	}
}

func TestJSONHasNoHeader(t *testing.T) {
	if languageOf(t, "package.json", "{}").HasHeader() {
		t.Error("JSON has a header")
	}
}

func languageOf(t *testing.T, path, content string) syntax.Language {
	t.Helper()
	language, found := syntax.LanguageOf(filepath.Base(path), content)
	if !found {
		t.Fatalf("no language for %s", path)
	}
	return language
}
