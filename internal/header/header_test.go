// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package header

import (
	"slices"
	"strings"
	"testing"

	"github.com/thehale/lawyer/internal/copyright"
	_ "github.com/thehale/lawyer/internal/languages"
	"github.com/thehale/lawyer/internal/syntax"
	"github.com/thehale/lawyer/internal/years"
)

var expectation = Expectation{Copyright: copyright.Expectation{Owner: "Joseph Hale", Years: years.InThePast(years.Only(2026))}, License: "MPL-2.0"}

func bash(t *testing.T) syntax.Language {
	language, found := syntax.LanguageOf("script.sh", "")
	if !found {
		t.Fatal("no language for script.sh")
	}
	return language
}

func TestCheck(t *testing.T) {
	cases := []struct {
		name, content string
		expectation   Expectation
		violations    []string
	}{
		{
			name: "correct",
			content: `
#!/usr/bin/env bash
# Copyright (c) 2026 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
		},
		{
			name: "missing",
			content: `
#!/usr/bin/env bash
echo
`,
			expectation: expectation,
			violations:  []string{"missing header"},
		},
		{
			name: "other licence",
			content: `
# Copyright (c) 2026 Joseph Hale
# SPDX-License-Identifier: MIT
`,
			expectation: expectation,
			violations:  []string{"expected MPL-2.0, found MIT"},
		},
		{
			name: "other owner",
			content: `
# Copyright (c) 2026 Someone Else
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
			violations:  []string{"missing copyright line for Joseph Hale"},
		},
		{
			name: "future",
			content: `
# Copyright (c) 2999 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
			violations:  []string{"years 2999 are in the future"},
		},
		{
			name: "unreadable",
			content: `
# Copyright (c) soon Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
			violations:  []string{`unreadable years in "soon Joseph Hale"`},
		},
		{
			name: "any past year",
			content: `
# Copyright (c) 2020 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
		},
		{
			name: "exact years",
			content: `
# Copyright (c) 2020 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: withExactYears(expectation),
			violations:  []string{"years 2020, expected 2026"},
		},
		{
			name: "copywrite's shape",
			content: `
# Copyright (c) Joseph Hale, 2026
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
			violations:  []string{`copyright line should read "Copyright (c) 2026 Joseph Hale"`},
		},
		{
			name: "spaced range",
			content: `
# Copyright (c) 2022 - 2025 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
			violations:  []string{`copyright line should read "Copyright (c) 2022-2025 Joseph Hale"`},
		},
		{
			name: "another license's notice",
			content: `
# Copyright (c) 2026 Joseph Hale
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
`,
			expectation: expectation,
			violations:  []string{"expected MPL-2.0, found the Apache-2.0 notice"},
		},
		{
			name: "its own notice",
			content: `
# Copyright (c) 2026 Joseph Hale
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
`,
			expectation: expectation,
			violations:  []string{"MPL-2.0 notice in place of an SPDX-License-Identifier line"},
		},
		{
			name: "prose",
			content: `
# Not a header: Copyright (c) 2026 Joseph Hale
# Not a header: SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
			violations:  []string{"missing header"},
		},
		{
			name: "below code",
			content: `
echo
# Copyright (c) 2026 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`,
			expectation: expectation,
			violations:  []string{"missing header"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if violations := Read(bash(t), text(c.content)).Violations(c.expectation); !slices.Equal(violations, c.violations) {
				t.Errorf("got %q, expected %q", violations, c.violations)
			}
		})
	}
}

func TestFixWritesExactYears(t *testing.T) {
	existing := Read(bash(t), text(`
# Copyright (c) 2020 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`))
	content := existing.ContentWith(existing.Canonical(withExactYears(expectation)))
	after := text(`
# Copyright (c) 2026 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`)

	if content != after {
		t.Errorf("got\n%s\nwant\n%s", content, after)
	}
}

func withExactYears(header Expectation) Expectation {
	header.Copyright.Years = years.Exactly(years.Only(2026))
	return header
}

func text(example string) string {
	return strings.TrimPrefix(example, "\n")
}
