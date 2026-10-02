// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"errors"
	"strings"

	"github.com/thehale/lawyer/internal/copyright"
	"github.com/thehale/lawyer/internal/header"
	"github.com/thehale/lawyer/internal/spdx"
	"github.com/thehale/lawyer/internal/years"
)

// A Declaration is what every Licensable's header must say: who owns the
// copyright, and under which license.
type Declaration struct {
	// Owner is the copyright holder as it should appear, on one line, such as
	// "Joseph Hale".
	Owner   string
	License Expression
}

// An Expression is an SPDX license expression, such as "MPL-2.0" or
// "MIT OR Apache-2.0".
type Expression = spdx.Expression

// Validate reports what is wrong with the Declaration, or nil when nothing is.
func (d Declaration) Validate() error {
	return errors.Join(d.ownerProblem(), d.licenseProblem())
}

// HeaderFor is the header d calls for in licensable, keeping the years its
// current header gives where d allows them.
func (d Declaration) HeaderFor(licensable Licensable) Header {
	return licensable.header.Canonical(d.headerExpectation(licensable))
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

func (d Declaration) licenseProblem() error {
	if d.License == "" {
		return errors.New("a license is required")
	} else {
		return nil
	}
}

func (d Declaration) headerExpectation(Licensable) header.Expectation {
	return header.Expectation{Copyright: copyright.Expectation{Owner: d.Owner, Years: years.InThePast(years.This())}, License: d.License}
}
