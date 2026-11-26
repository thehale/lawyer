// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"slices"
	"strings"
)

var argumentOptions = []string{"--copyright-owner", "--license", "--copyright-year", "--exclude"}

func nextOption(args []string) (flag, argument string, rest []string, err error) {
	flag, argument, inline := strings.Cut(args[0], "=")

	switch {
	case flag == "-" || !strings.HasPrefix(flag, "-"):
		return "", args[0], args[1:], nil
	case inline && flag == "--fix":
		return flag, "", nil, fmt.Errorf("%s takes no value", flag)
	case inline || !slices.Contains(argumentOptions, flag):
		return flag, argument, args[1:], nil
	case len(args) > 1:
		return flag, args[1], args[2:], nil
	default:
		return flag, "", nil, fmt.Errorf("%s needs a value", flag)
	}
}

func set[T ~string](option *T, flag, argument string) error {
	if *option == "" {
		*option = T(argument)
		return nil
	} else {
		return fmt.Errorf("%s given twice", flag)
	}
}
