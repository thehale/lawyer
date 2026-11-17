// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"cmp"
	"errors"
	"fmt"

	"github.com/thehale/lawyer"
)

func validate(declaration lawyer.Declaration) error {
	owner := require("--copyright-owner", declaration.Owner, "names an owner")
	license := require("--license", string(declaration.License), "holds a license lawyer knows")
	return cmp.Or(errors.Join(owner, license), declaration.Validate())
}

func require(flag, value, fallback string) error {
	if value == "" {
		return fmt.Errorf("%s is required, as no LICENSE file %s", flag, fallback)
	} else {
		return nil
	}
}
