// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package paths

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

func FilesWithin(root Path, files []Path) ([]Path, error) {
	base, err := realPath(string(root))
	here, hereErr := realPath(".")
	var within []Path
	var outside []error
	for _, path := range files {
		if isInside(base, here, string(path)) {
			within = append(within, path)
		} else {
			outside = append(outside, fmt.Errorf("%s: outside the repository", path.Printable()))
		}
	}
	return within, errors.Join(append(outside, err, hereErr)...)
}

func isInside(base, here, path string) bool {
	relative, err := filepath.Rel(base, location(here, path))
	escapes := relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator))
	return err == nil && !escapes
}

func location(here, path string) string {
	parent, err := realPath(filepath.Dir(path))

	if err == nil {
		return filepath.Join(parent, filepath.Base(path))
	} else {
		return filepath.Join(here, lexicalPath(path))
	}
}

func lexicalPath(path string) string {
	here, err := filepath.Abs(".")
	absolute, absErr := filepath.Abs(path)
	relative, relErr := filepath.Rel(here, absolute)

	if err == nil && absErr == nil && relErr == nil {
		return relative
	} else {
		return ".."
	}
}

func realPath(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err == nil {
		absolute, err = filepath.EvalSymlinks(absolute)
	}
	return absolute, err
}
