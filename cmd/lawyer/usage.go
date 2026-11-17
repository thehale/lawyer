// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
	"slices"
)

const usage = `Usage: lawyer check [--fix] [--copyright-owner OWNER] [--license EXPRESSION]
                    [--copyright-year YEARS] [--exclude GLOB]... [--] [PATH... | -]
       lawyer languages
       lawyer licenses

check      Checks that files carry a copyright and SPDX header and that the
           LICENSE matches, and --fix writes them. The owner and license
           default to what the LICENSE files name.
languages  Lists the languages lawyer puts headers in.
licenses   Lists the SPDX ids --license accepts.

Examples:
  lawyer check --copyright-owner "Your Name" --license MPL-2.0 --fix
  lawyer check --exclude '**/vendor/**'
  git diff --name-only --diff-filter=d main | lawyer check -
`

func isHelp(arg string) bool {
	return slices.Contains([]string{"-h", "--help"}, arg)
}

func printUsage() int {
	fmt.Print(usage)
	return 0
}

func misuse(err error) int {
	fmt.Fprintf(os.Stderr, "lawyer: %v\n\n%s", err, usage)
	return 2
}
