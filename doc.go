// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

// Package lawyer checks and fixes the copyright and SPDX license headers of a
// repository's files.
//
// A Declaration says who owns the copyright, under which license, and for
// which years (by default, each file's years from its git history). A
// Repository holds the Licensables whose headers must carry it:
//
//	declaration := lawyer.Declaration{Owner: "Joseph Hale", License: "MPL-2.0", Years: "2024-2026"}
//	repository := lawyer.NewRepository(".")
//
//	violations, err := repository.Check(declaration)
//	repository.Fix(declaration)
//
// If you want more granular insight into the violations...
//
//	headerViolations, err := repository.Licensables().Excluding("**/vendor/**").Check(declaration)
//
// If you want more granular control over fixes...
//
//	headerChanges, err := repository.Licensables().Excluding("**/vendor/**").Fix(declaration)
package lawyer
