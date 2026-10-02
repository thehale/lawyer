// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

//go:build !unix

package files

import (
	"io/fs"
	"os"
)

func own(*os.File, fs.FileInfo) error {
	return nil
}
