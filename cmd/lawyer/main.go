// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
)

var version = "dev"

const usage = `Usage: lawyer check [--fix] --copyright-owner OWNER --license EXPRESSION
                    [--copyright-year YEARS] [--exclude GLOB]... [PATH... | -]
       lawyer languages
       lawyer licenses

check      Checks that files carry a copyright and SPDX header and that the
           LICENSE matches, and --fix writes them.
languages  Lists the languages lawyer puts headers in.
licenses   Lists the SPDX ids --license accepts.
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	switch append(args, "")[0] {
	case "check":
		return runCheck(args[1:])
	case "languages":
		return listLanguages()
	case "licenses":
		return listLicenses()
	case "--help", "-h", "help":
		fmt.Print(usage)
		return 0
	case "--version":
		fmt.Println(version)
		return 0
	default:
		fmt.Fprint(os.Stderr, usage)
		return 2
	}
}
