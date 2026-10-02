// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/thehale/lawyer/internal/copyright"
	"github.com/thehale/lawyer/internal/spdx"
	"github.com/thehale/lawyer/internal/years"
)

var expectation = copyright.Expectation{Owner: "Joseph Hale", Years: years.InThePast(years.Only(2026))}

func byID(t *testing.T, id spdx.ID) spdx.License {
	license, found := spdx.ByID(id)
	if !found {
		t.Fatalf("no license %s", id)
	}
	return license
}

func TestRenderedLicensesPass(t *testing.T) {
	for _, id := range spdx.IDs() {
		license := byID(t, id)
		if violations := file(file("").canonical(license, expectation)).violations(license, expectation); violations != nil {
			t.Errorf("%s: %q", id, violations)
		}
	}
}

func TestTheTemplatesLicensePasses(t *testing.T) {
	content, err := os.ReadFile("LICENSE")
	if err != nil {
		t.Fatal(err)
	}
	if violations := file(string(content)).violations(byID(t, "MPL-2.0"), expectation); violations != nil {
		t.Errorf("%q", violations)
	}
}

func TestMITPutsTheCopyrightInItsOwnPlace(t *testing.T) {
	licenseText := file("").canonical(byID(t, "MIT"), expectation)

	top := text(`
MIT License

Copyright (c) 2026 Joseph Hale

Permission`)

	if !strings.HasPrefix(licenseText, top) {
		t.Errorf("got %q", licenseText[:80])
	}
}

func TestFormattingIsTolerated(t *testing.T) {
	license := byID(t, "MIT")
	licenseText := file("").canonical(license, expectation)
	markdown := "# " + strings.ReplaceAll(strings.ReplaceAll(licenseText, "\n\n", "\n\n  "), `"`, "“")

	if violations := file(markdown).violations(license, expectation); violations != nil {
		t.Errorf("%q", violations)
	}
}

func TestViolations(t *testing.T) {
	license := byID(t, "MIT")
	licenseText := file("").canonical(license, expectation)
	cases := []struct {
		name, content string
		expectation   copyright.Expectation
		violations    []string
	}{
		{"other owner", strings.Replace(licenseText, "Joseph Hale", "Someone Else", 1), expectation, []string{"missing copyright line for Joseph Hale"}},
		{"no copyright", license.Text, expectation, []string{"missing copyright line for Joseph Hale"}},
		{"reworded", strings.Replace(licenseText, "without restriction", "with restrictions", 1), expectation, []string{"text differs from MIT"}},
		{"another license", file("").canonical(byID(t, "MPL-2.0"), expectation), expectation, []string{"text differs from MIT"}},
		{"future", strings.Replace(licenseText, "2026", "2999", 1), expectation, []string{"years 2999 are in the future"}},
		{"exact", licenseText, copyright.Expectation{Owner: "Joseph Hale", Years: years.Exactly(years.Range{First: 2024, Last: 2026})}, []string{"years 2026, expected 2024-2026"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if violations := file(c.content).violations(license, c.expectation); !slices.Equal(violations, c.violations) {
				t.Errorf("got %q, expected %q", violations, c.violations)
			}
		})
	}
}

func TestRenderKeepsPastYears(t *testing.T) {
	existing := text(`
Copyright (c) 2020-2021 Joseph Hale

old text`)
	top := text(`
MIT License

Copyright (c) 2020-2021 Joseph Hale`)
	licenseText := file(existing).canonical(byID(t, "MIT"), expectation)

	if !strings.HasPrefix(licenseText, top) {
		t.Errorf("got %q", licenseText[:60])
	}
}

func TestLicensesWithoutASlotTakeNoCopyrightLine(t *testing.T) {
	for _, id := range []spdx.ID{"MPL-2.0", "Apache-2.0"} {
		license := byID(t, id)
		added := "Copyright (c) 2026 Joseph Hale\n\n" + license.Text

		if file("").canonical(license, expectation) != license.Text || file(license.Text).violations(license, expectation) != nil || file(added).violations(license, expectation) == nil {
			t.Errorf("%s: a copyright line is rendered or required", id)
		}
	}
}

func text(example string) string {
	return strings.TrimPrefix(example, "\n")
}

func file(content string) licenseFile {
	return licenseFile{content: content}
}
