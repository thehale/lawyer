// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"slices"
	"strings"
)

var valueOptions = []string{"--copyright-owner", "--license"}

func nextOption(args []string) (flag, value string, rest []string, err error) {
	flag, value, inline := strings.Cut(args[0], "=")

	switch {
	case !strings.HasPrefix(flag, "--"):
		return "", args[0], args[1:], nil
	case inline || !slices.Contains(valueOptions, flag):
		return flag, value, args[1:], nil
	case len(args) > 1:
		return flag, args[1], args[2:], nil
	default:
		return flag, "", nil, fmt.Errorf("%s needs a value", flag)
	}
}

func set[T ~string](option *T, flag, value string) error {
	if *option == "" {
		*option = T(value)
		return nil
	} else {
		return fmt.Errorf("%s given twice", flag)
	}
}
