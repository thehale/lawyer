// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"fmt"
	"slices"
	"testing"

	"github.com/thehale/lawyer/internal/spdx"
)

func TestLicensingViolations(t *testing.T) {
	mpl, mit := licenseText(t, "MPL-2.0"), licenseText(t, "MIT")
	cases := []struct {
		name       string
		license    Expression
		files      []licenseFile
		violations []string
	}{
		{"canonical", "MPL-2.0", []licenseFile{{"LICENSE", mpl, "LICENSE", true}}, nil},
		{"any case and extension", "MPL-2.0", []licenseFile{{"license.md", mpl, "license.md", true}}, nil},
		{"missing", "MPL-2.0", nil, []string{"LICENSE: missing"}},
		{"another license's text", "MPL-2.0", []licenseFile{{"LICENSE", mit, "LICENSE", true}}, []string{"LICENSE: text differs from MPL-2.0"}},
		{"several variants", "MPL-2.0", []licenseFile{{"LICENSE", mpl, "LICENSE", true}, {"LICENSE.md", mpl, "LICENSE.md", true}}, []string{"LICENSE: one of several files for MPL-2.0", "LICENSE.md: one of several files for MPL-2.0"}},
		{"not a regular file", "MPL-2.0", []licenseFile{{"LICENSE", "", "LICENSE", false}}, []string{"LICENSE: not a regular file"}},
		{"a stray", "MPL-2.0", []licenseFile{{"LICENSE", mpl, "LICENSE", true}, {"LICENSE-MIT.md", mit, "LICENSE-MIT.md", true}}, []string{"LICENSE-MIT.md: " + uncalledFor}},
		{"several licenses replace the plain one", "MIT AND Apache-2.0", []licenseFile{{"LICENSE", mpl, "LICENSE", true}}, []string{"LICENSE-MIT: missing", "LICENSE-Apache-2.0: missing", "LICENSE: " + uncalledFor}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			licensing := Licensing{repository: Repository{root: "."}, files: c.files}
			declaration := Declaration{Owner: "Joseph Hale", License: c.license, Years: "2026"}
			violations, _ := licensing.Check(declaration)
			if messages := messages(violations); !slices.Equal(messages, c.violations) {
				t.Errorf("got %q, expected %q", messages, c.violations)
			}
		})
	}
}

func messages(violations []Violation) []string {
	var messages []string
	for _, violation := range violations {
		messages = append(messages, fmt.Sprintf("%s: %s", violation.Path, violation.Message))
	}
	return messages
}

func licenseText(t *testing.T, id spdx.ID) string {
	license, found := spdx.ByID(id)
	if !found {
		t.Fatalf("no license %s", id)
	}
	return license.Text
}
