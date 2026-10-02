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
	options, err := parse(args)

	if err == nil {
		return options.run(lawyer.NewRepository("."))
	} else {
		fmt.Fprintf(os.Stderr, "lawyer: %v\n\n%s", err, usage)
		return 2
	}
}

type options struct {
	fix         bool
	declaration lawyer.Declaration
	paths       []lawyer.Path
}

func parse(args []string) (options, error) {
	var chosen options
	var err error
	for rest := args; len(rest) > 0 && err == nil; {
		rest, err = chosen.take(rest)
	}
	return chosen, errors.Join(err, chosen.declaration.Validate())
}

func (o *options) take(args []string) ([]string, error) {
	flag, value, rest, err := nextOption(args)

	switch flag {
	case "":
		o.paths = append(o.paths, lawyer.Path(value))
	case "--fix":
		o.fix = true
	case "--copyright-owner":
		err = errors.Join(err, set(&o.declaration.Owner, flag, value))
	case "--license":
		err = errors.Join(err, set(&o.declaration.License, flag, value))
	default:
		err = fmt.Errorf("unknown option %s", flag)
	}
	return rest, err
}

func (o options) run(repository lawyer.Repository) int {
	licensables := repository.Licensables(o.paths...)

	if o.fix {
		changes, err := licensables.Fix(o.declaration)
		return printChanges(changes, err)
	} else {
		violations, err := licensables.Check(o.declaration)
		return printViolations(violations, err)
	}
}
