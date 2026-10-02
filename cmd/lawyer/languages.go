// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/thehale/lawyer"
)

func listLanguages() int {
	table := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	var failures []error
	for _, language := range lawyer.Languages() {
		_, err := fmt.Fprintf(table, "%s\t%s\n", language.Name, selectors(language))
		failures = append(failures, err)
	}
	return printError(errors.Join(append(failures, table.Flush())...))
}

func selectors(language lawyer.Language) string {
	names := slices.Clone(language.Filenames)
	for _, extension := range language.Extensions {
		names = append(names, "."+extension)
	}
	for _, interpreter := range language.Shebangs {
		names = append(names, "#!"+interpreter)
	}
	if !language.HasHeader {
		names = append(names, "(no header)")
	}
	return strings.Join(names, " ")
}
