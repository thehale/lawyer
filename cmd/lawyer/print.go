// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	"github.com/thehale/lawyer"
)

func printViolations(violations []lawyer.Violation, err error) int {
	for _, violation := range violations {
		fmt.Fprintf(os.Stderr, "%s: %s\n", violation.Path, violation.Message)
	}
	code := printError(err)

	if len(violations) > 0 {
		return 1
	} else {
		return code
	}
}

func printError(err error) int {
	if err != nil {
		fmt.Fprintf(os.Stderr, "lawyer: %v\n", err)
		return 1
	} else {
		return 0
	}
}
