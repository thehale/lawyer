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
	command, rest := append(args, "")[0], args[min(1, len(args)):]

	switch command {
	case "check":
		return runCheck(rest)
	case "languages":
		return withoutArguments(rest, listLanguages)
	case "licenses":
		return withoutArguments(rest, listLicenses)
	case "", "help", "-h", "--help":
		return printUsage()
	case "--version":
		fmt.Println(version)
		return 0
	default:
		fmt.Fprint(os.Stderr, usage)
		return 2
	}
}

func withoutArguments(args []string, list func() int) int {
	if len(args) == 1 && isHelp(args[0]) {
		return printUsage()
	} else {
		return list()
	}
}
