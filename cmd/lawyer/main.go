// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
)

var version = "dev"

const usage = `Usage: lawyer check [--fix] --copyright-owner OWNER --license EXPRESSION
                    [--exclude GLOB]... [PATH... | -]

Checks that files carry a copyright and SPDX header, and --fix writes them.
`

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	switch append(args, "")[0] {
	case "check":
		return runCheck(args[1:])
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
