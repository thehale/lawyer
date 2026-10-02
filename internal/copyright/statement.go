// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package copyright

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/thehale/lawyer/internal/years"
)

// A Statement is what a copyright line says after its "Copyright", "(c)" or
// "©": its years and holders, in any common order and punctuation.
type Statement string

var (
	copyrightLine  = regexp.MustCompile(`^(?i:copyright\b(?:\s*(?:\(c\)|©))?|\(c\)|©)\s*(.*)$`)
	rightsReserved = regexp.MustCompile(`(?i)all rights reserved`)
	yearRun        = regexp.MustCompile(`\b[0-9]{4}(?:\s*[-–,]\s*[0-9]{4})*\b`)
)

func StatementIn(line string) (Statement, bool) {
	match := copyrightLine.FindStringSubmatch(strings.TrimSpace(line))

	if match == nil {
		return "", false
	} else {
		return Statement(match[1]), true
	}
}

func StatementBy(text, owner string) (Statement, bool) {
	for _, line := range strings.Split(text, "\n") {
		if statement, ok := StatementIn(line); ok && statement.HasOwner(owner) {
			return statement, true
		}
	}
	return "", false
}

func (s Statement) HasOwner(owner string) bool {
	return strings.Contains(strings.ToLower(s.Owner()), strings.ToLower(owner))
}

func (s Statement) HasYear() bool {
	_, err := s.Years()
	return err == nil
}

// Years are the years the statement gives, such as 2022-2024 for "2022, 2024
// Jane Doe. All rights reserved.".
func (s Statement) Years() (years.Range, error) {
	written := yearRun.FindString(string(s))

	if written != "" {
		return years.LooseRange(written)
	} else {
		return years.Range{}, fmt.Errorf("unreadable years in %q", string(s))
	}
}

// Owner is who the statement names, such as "Jane Doe" for "2022, 2024 Jane
// Doe. All rights reserved.".
func (s Statement) Owner() string {
	rest := rightsReserved.ReplaceAllString(yearRun.ReplaceAllString(string(s), ""), "")
	return strings.Join(strings.Fields(strings.Trim(rest, " ,.;")), " ")
}

// Canonical is the line expectation calls for in the statement's place,
// keeping its years where expectation allows them.
func (s Statement) Canonical(expectation Expectation) Line {
	current, readErr := s.Years()
	return Line{Years: expectation.Years.Replacement(current, readErr), Owner: expectation.Owner}
}
