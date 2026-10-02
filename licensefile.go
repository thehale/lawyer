// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/thehale/lawyer/internal/copyright"
	"github.com/thehale/lawyer/internal/spdx"
	"github.com/thehale/lawyer/internal/years"
)

// A licenseFile is one of a repository's LICENSE files.
type licenseFile struct {
	name, content string
	path          Path
	regular       bool
}

var (
	recognizedName = regexp.MustCompile(`(?i)^(LICEN[CS]E|COPYING)(-[a-z0-9.+-]+)?(\.(md|txt))?$`)
	managedName    = regexp.MustCompile(`(?i)^LICEN[CS]E(-[a-z0-9.+-]+)?(\.(md|txt))?$`)
)

func readLicenseFile(root Path, name string) (licenseFile, error) {
	path := root.Child(name)
	info, err := os.Lstat(string(path))

	if err == nil && info.Mode().IsRegular() {
		content, err := os.ReadFile(string(path))
		return licenseFile{name: name, path: path, content: string(content), regular: true}, err
	} else {
		return licenseFile{name: name, path: path}, err
	}
}

// license is the license the file holds: the one its name gives, as
// LICENSE-MIT does, or else the one its text is.
// violations are how the file falls short of holding license as expectation
// calls for.
func (f licenseFile) violations(license spdx.License, expectation copyright.Expectation) []string {
	written, owned := license.OwnerYears(f.content, expectation.Owner)
	slotted := license.HasCopyrightSlot()

	switch {
	case slotted && owned:
		return expectation.YearViolations(years.RangeOf(written))
	case slotted && license.HasCopyrightLine(f.content):
		return []string{fmt.Sprintf("missing copyright line for %s", expectation.Owner)}
	case license.IsWrittenIn(f.content):
		return nil
	default:
		return []string{fmt.Sprintf("text differs from %s", license.ID)}
	}
}

// canonical is the text the file should hold: license's, with the copyright
// line expectation calls for, keeping its years where expectation allows.
func (f licenseFile) canonical(license spdx.License, expectation copyright.Expectation) string {
	statement, _ := copyright.StatementBy(f.content, expectation.Owner)
	return license.TextWithCopyright(statement.Canonical(expectation).String())
}

// isLicenseName reports whether name is a license file's name, such as
// LICENSE, COPYING or LICENSE-MIT.md. A license file takes no header.
func isLicenseName(name string) bool {
	return recognizedName.MatchString(name)
}

// isManagedName reports whether name is one lawyer writes license files
// under: LICENSE or LICENSE-<id>, plain, .md or .txt.
func isManagedName(name string) bool {
	return managedName.MatchString(name)
}

// licenseNames are the names license's file may have among an expression's
// licenses: LICENSE alone, or LICENSE-<id> beside others, plain, .md or .txt.
func licenseNames(license spdx.License, among []spdx.License) []string {
	base := "LICENSE"
	if len(among) > 1 {
		base += "-" + string(license.ID)
	}
	return []string{base, base + ".md", base + ".txt"}
}

func hasName(file licenseFile, names []string) bool {
	return slices.ContainsFunc(names, func(name string) bool { return strings.EqualFold(name, file.name) })
}
