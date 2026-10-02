// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/thehale/lawyer/internal/copyright"
	"github.com/thehale/lawyer/internal/header"
	"github.com/thehale/lawyer/internal/history"
	"github.com/thehale/lawyer/internal/spdx"
	"github.com/thehale/lawyer/internal/years"
)

// A Declaration is what every Licensable's header and the repository's
// Licensing must say: who owns the copyright, under which license, and for
// which years.
type Declaration struct {
	// Owner is the copyright holder as it should appear, on one line, such as
	// "Joseph Hale".
	Owner   string
	License Expression
	// Years is the year or range every header must state, such as
	// "2024-2026". Left empty, each file's years come from its git history.
	Years string
}

// An Expression is an SPDX license expression, such as "MPL-2.0" or
// "MIT OR Apache-2.0".
type Expression = spdx.Expression

// Validate reports what is wrong with the Declaration, or nil when nothing is.
func (d Declaration) Validate() error {
	_, licenseErr := d.License.Licenses()
	_, yearsErr := d.span()
	return errors.Join(d.ownerProblem(), d.licenseProblem(licenseErr), yearsErr)
}

// HeaderFor is the header d calls for in licensable, keeping the years its
// current header gives where d allows them.
func (d Declaration) HeaderFor(licensable Licensable) Header {
	return licensable.header.Canonical(d.headerExpectation(licensable))
}

// LicensingFor is the Licensing d calls for in repository, keeping its LICENSE
// files' names and past years where d allows them.
func (d Declaration) LicensingFor(repository Repository) Licensing {
	return repository.Licensing().declaredBy(d)
}

func (d Declaration) ownerProblem() error {
	switch {
	case d.Owner == "":
		return errors.New("an owner is required")
	case strings.ContainsAny(d.Owner, "\r\n"):
		return errors.New("the owner must fit on one line")
	default:
		return nil
	}
}

func (d Declaration) licenseProblem(licenseErr error) error {
	switch {
	case d.License == "":
		return errors.New("a license is required")
	case licenseErr != nil:
		return fmt.Errorf("license: %w", licenseErr)
	default:
		return nil
	}
}

func (d Declaration) span() (years.Range, error) {
	span, err := years.RangeOf(d.Years)

	switch {
	case d.Years == "":
		return years.Range{}, nil
	case err != nil:
		return years.Range{}, fmt.Errorf("years: %w", err)
	default:
		return span, nil
	}
}

func (d Declaration) headerExpectation(licensable Licensable) header.Expectation {
	history := licensable.repository.history
	return header.Expectation{Copyright: d.copyright(history, history.EditYears(licensable.path)), License: d.License}
}

func (d Declaration) licenseExpectation(repository Repository) copyright.Expectation {
	return d.copyright(repository.history, repository.history.Years())
}

func (d Declaration) copyright(history history.History, editYears years.Range) copyright.Expectation {
	span, _ := d.span()

	switch {
	case d.Years != "":
		return copyright.Expectation{Owner: d.Owner, Years: years.Exactly(span)}
	case history.IsEveryEditDated():
		return copyright.Expectation{Owner: d.Owner, Years: years.Exactly(editYears)}
	default:
		return copyright.Expectation{Owner: d.Owner, Years: years.InThePast(editYears)}
	}
}
