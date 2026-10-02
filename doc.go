// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

// Package lawyer checks and fixes the copyright and SPDX license headers of a
// repository's files.
//
// A Declaration says who owns the copyright and under which license. A
// Repository holds the Licensables whose headers must carry it:
//
//	declaration := lawyer.Declaration{Owner: "Joseph Hale", License: "MPL-2.0"}
//	repository := lawyer.NewRepository(".")
//
//	violations, err := repository.Check(declaration)
//	repository.Fix(declaration)
//
// If you want to check or fix only some of the files...
//
//	headerViolations, err := repository.Licensables("bin", "install.sh").Check(declaration)
//	headerChanges, err := repository.Licensables("bin", "install.sh").Fix(declaration)
package lawyer
