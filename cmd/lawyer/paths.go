// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/thehale/lawyer"
)

func paths(args []string, stdin io.Reader) ([]lawyer.Path, error) {
	var files []lawyer.Path
	var failures []error
	for _, arg := range args {
		named, err := pathsOf(arg, stdin)
		files = append(files, named...)
		failures = append(failures, err)
	}
	return files, errors.Join(failures...)
}

func pathsOf(arg string, stdin io.Reader) ([]lawyer.Path, error) {
	if arg == "-" {
		return lines(stdin)
	} else {
		return []lawyer.Path{lawyer.Path(arg)}, nil
	}
}

func lines(stdin io.Reader) ([]lawyer.Path, error) {
	var files []lawyer.Path
	scanner := bufio.NewScanner(stdin)
	for scanner.Scan() {
		if line := strings.TrimSuffix(scanner.Text(), "\r"); line != "" {
			files = append(files, lawyer.Path(line))
		}
	}
	return files, scanner.Err()
}
