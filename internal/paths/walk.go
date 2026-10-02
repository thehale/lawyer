// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package paths

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/thehale/lawyer/internal/git"
)

func FilesUnder(root Path) ([]Path, error) {
	if git.IsInCheckout(string(root)) {
		return visibleFiles(root)
	} else {
		return allFiles(root)
	}
}

func visibleFiles(root Path) ([]Path, error) {
	output, err := git.Run(".", "ls-files", "-z", "--cached", "--others", "--exclude-standard", "--", string(root))
	var files []Path
	for _, file := range strings.Split(strings.TrimSuffix(output, "\x00"), "\x00") {
		files = append(files, Path(file))
	}
	return slices.DeleteFunc(Distinct(files), isMissing), err
}

func isMissing(path Path) bool {
	_, err := os.Lstat(string(path))
	return path == "" || err != nil
}

func allFiles(root Path) ([]Path, error) {
	var files []Path
	err := filepath.WalkDir(string(root), func(path string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() {
			files = append(files, Path(path))
		}
		return err
	})
	return files, err
}
