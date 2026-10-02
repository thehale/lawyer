// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"slices"
	"testing"
)

func TestExpressionLicenses(t *testing.T) {
	cases := map[string][]ID{
		"MPL-2.0":                     {"MPL-2.0"},
		"MIT OR Apache-2.0":           {"MIT", "Apache-2.0"},
		"(MIT AND Apache-2.0) OR MIT": {"MIT", "Apache-2.0"},
	}
	for text, ids := range cases {
		licenses, err := Expression(text).Licenses()
		if err != nil || !slices.Equal(idsOf(licenses), ids) {
			t.Errorf("%s: got %q, %v", text, idsOf(licenses), err)
		}
	}
}

func TestRejectedExpressions(t *testing.T) {
	for _, expression := range []string{"", "MIT OR", "OR MIT", "(MIT", "MIT)", "MIT MIT", "CC-BY-4.0", "GPL-2.0-only WITH Classpath-exception-2.0", "LicenseRef-mine"} {
		if _, err := Expression(expression).Licenses(); err == nil {
			t.Errorf("%q accepted", expression)
		}
	}
}

func idsOf(found []License) []ID {
	var ids []ID
	for _, license := range found {
		ids = append(ids, license.ID)
	}
	return ids
}
