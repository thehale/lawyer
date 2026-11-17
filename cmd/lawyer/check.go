// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/thehale/lawyer"
)

func runCheck(args []string) int {
	repository := lawyer.NewRepository(".")
	options, err := parse(args, repository.Licensing())

	switch {
	case options.help:
		return printUsage()
	case err == nil:
		return options.run(repository)
	default:
		return misuse(err)
	}
}

type options struct {
	help        bool
	fix         bool
	declaration lawyer.Declaration
	excludes    []lawyer.Glob
	paths       []string
}

func parse(args []string, licensing lawyer.Licensing) (options, error) {
	var chosen options
	var err error
	for rest := args; len(rest) > 0 && err == nil; {
		rest, err = chosen.take(rest)
	}
	if chosen.declaration.Owner == "" || chosen.declaration.License == "" {
		chosen.declaration = chosen.declaration.Or(lawyer.DeclarationFrom(licensing))
	}
	return chosen, errors.Join(err, validate(chosen.declaration))
}

func (o *options) take(args []string) ([]string, error) {
	flag, value, rest, err := nextOption(args)

	switch flag {
	case "":
		o.paths = append(o.paths, value)
	case "--":
		o.paths, rest = append(o.paths, rest...), nil
	case "-h", "--help":
		o.help = true
	case "--fix":
		o.fix = true
	case "--exclude":
		o.excludes = append(o.excludes, lawyer.Glob(value))
	case "--copyright-owner":
		err = errors.Join(err, set(&o.declaration.Owner, flag, value))
	case "--license":
		err = errors.Join(err, set(&o.declaration.License, flag, value))
	case "--copyright-year":
		err = errors.Join(err, set(&o.declaration.Years, flag, value))
	default:
		err = fmt.Errorf("unknown option %s", flag)
	}
	return rest, err
}

func (o options) run(repository lawyer.Repository) int {
	licensables, pathsErr := o.licensables(repository)
	licensing := repository.Licensing()
	printWarnings(repository.Warnings(o.declaration))

	if o.fix {
		headers, headersErr := licensables.Fix(o.declaration)
		licenses, licensesErr := licensing.Fix(o.declaration)
		return printChanges(headers.Union(licenses), errors.Join(pathsErr, headersErr, licensesErr))
	} else {
		headers, headersErr := licensables.Check(o.declaration)
		licenses, licensesErr := licensing.Check(o.declaration)
		return printViolations(append(headers, licenses...), errors.Join(pathsErr, headersErr, licensesErr))
	}
}

func (o options) licensables(repository lawyer.Repository) (lawyer.Licensables, error) {
	if len(o.paths) == 0 {
		return repository.Licensables().Excluding(o.excludes...), nil
	} else {
		targets, pathsErr := paths(o.paths, os.Stdin)
		return repository.Licensables(append(targets, "")...).Excluding(o.excludes...), pathsErr
	}
}
