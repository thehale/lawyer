// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package copyright

import "testing"

func TestStatementOwner(t *testing.T) {
	cases := map[string]string{
		"Copyright (c) 2018 Joseph Hale":              "Joseph Hale",
		"Copyright (c) 2022 - 2025 Joseph Hale":       "Joseph Hale",
		"© 2022, 2024 Jane Doe. All rights reserved.": "Jane Doe",
		"Copyright 2019-2021, The Go Authors":         "The Go Authors",
	}
	for line, owner := range cases {
		statement, _ := StatementIn(line)
		if got := statement.Owner(); got != owner {
			t.Errorf("%q: got %q, expected %q", line, got, owner)
		}
	}
}
