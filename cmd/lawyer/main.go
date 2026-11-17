// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
)

var version = "dev"

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
		return printUsage()
	case "--version":
		fmt.Println(version)
		return 0
	default:
		fmt.Fprint(os.Stderr, usage)
		return 2
	}
}
