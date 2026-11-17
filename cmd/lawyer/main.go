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
		return withoutArguments(command, rest, listLanguages)
	case "licenses":
		return withoutArguments(command, rest, listLicenses)
	case "", "help", "-h", "--help":
		return printUsage()
	case "version", "--version":
		fmt.Println(version)
		return 0
	default:
		return misuse(fmt.Errorf("unknown command %s", command))
	}
}

func withoutArguments(command string, args []string, list func() int) int {
	switch {
	case len(args) == 0:
		return list()
	case len(args) == 1 && isHelp(args[0]):
		return printUsage()
	default:
		return misuse(fmt.Errorf("%s takes no arguments", command))
	}
}
