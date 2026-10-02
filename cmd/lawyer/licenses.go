// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"

	"github.com/thehale/lawyer"
)

func listLicenses() int {
	for _, id := range lawyer.Licenses() {
		fmt.Println(id)
	}
	return 0
}
