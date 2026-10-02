// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	spacing    = strings.NewReplacer("(", " ( ", ")", " ) ")
	wellFormed = regexp.MustCompile(`^[(\s]*[\w.+-]+(?:[()\s]+(?:AND|OR|WITH)[(\s]+[\w.+-]+)*[)\s]*$`)
)

// An Expression is an SPDX license expression, such as "MPL-2.0" or
// "MIT OR Apache-2.0".
type Expression string

// IsWellFormed reports whether the expression follows SPDX's grammar, whether or
// not lawyer carries the licenses it names.
func (e Expression) IsWellFormed() bool {
	return wellFormed.MatchString(string(e))
}

// Licenses are the licenses the expression names, each once, or an error when
// the expression is malformed or names a license lawyer doesn't carry.
func (e Expression) Licenses() ([]License, error) {
	tokens := strings.Fields(spacing.Replace(string(e)))
	err := grammarError(tokens)

	if err == nil {
		return knownLicenses(ids(tokens))
	} else {
		return nil, err
	}
}

func knownLicenses(ids []ID) ([]License, error) {
	var licenses []License
	var failures []error
	for _, id := range ids {
		found, ok := ByID(id)
		if !ok {
			failures = append(failures, fmt.Errorf("unknown license %s", id))
		}
		licenses = append(licenses, found)
	}
	return licenses, errors.Join(failures...)
}

func grammarError(tokens []string) error {
	depth := 0
	wantsID := true
	var err error
	for index := 0; index < len(tokens) && err == nil; index++ {
		token := tokens[index]
		switch {
		case token == "WITH":
			err = fmt.Errorf("WITH exceptions are not supported yet")
		case strings.HasPrefix(token, "LicenseRef-") || strings.HasPrefix(token, "DocumentRef-"):
			err = fmt.Errorf("%s: custom licenses are not supported yet", token)
		case token == "(" && wantsID:
			depth++
		case token == ")" && !wantsID && depth > 0:
			depth--
		case (token == "AND" || token == "OR") && !wantsID:
			wantsID = true
		case !isOperator(token) && wantsID:
			wantsID = false
		default:
			err = fmt.Errorf("unexpected %q in license expression", token)
		}
	}
	if err == nil && (wantsID || depth != 0) {
		err = fmt.Errorf("incomplete license expression %q", strings.Join(tokens, " "))
	}
	return err
}

func isOperator(token string) bool {
	return slices.Contains([]string{"(", ")", "AND", "OR", "WITH"}, token)
}

func ids(tokens []string) []ID {
	var unique []ID
	for _, token := range tokens {
		if !isOperator(token) && !slices.Contains(unique, ID(token)) {
			unique = append(unique, ID(token))
		}
	}
	return unique
}
