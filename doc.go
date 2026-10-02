// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

// Package lawyer checks the copyright and SPDX license headers of a
// repository's files.
//
// A Declaration says who owns the copyright and under which license. A
// Repository holds the Licensables whose headers must carry it:
//
//	declaration := lawyer.Declaration{Owner: "Joseph Hale", License: "MPL-2.0"}
//	repository := lawyer.NewRepository(".")
//
//	violations, err := repository.Licensables("install.sh").Check(declaration)
package lawyer
