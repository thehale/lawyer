// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/thehale/lawyer/internal/paths"
)

func Replace(path paths.Path, before, after string) error {
	temporary, err := os.CreateTemp(filepath.Dir(string(path)), ".lawyer-*")
	if err == nil {
		err = move(temporary, path, before, after)
		_ = os.Remove(temporary.Name())
	}
	return err
}

func move(temporary *os.File, path paths.Path, before, after string) error {
	original, _ := os.Lstat(string(path))
	err := write(temporary, after, original)
	if err == nil && !isUnchanged(path, before) {
		err = fmt.Errorf("%s changed while being fixed, so it was left alone", path.Printable())
	} else if err == nil {
		err = os.Rename(temporary.Name(), string(path))
	}
	return err
}

func isUnchanged(path paths.Path, before string) bool {
	now, err := os.ReadFile(string(path))
	return err == nil && string(now) == before
}

func write(temporary *os.File, content string, original fs.FileInfo) error {
	_, err := temporary.WriteString(content)
	return errors.Join(err, own(temporary, original), temporary.Chmod(mode(original)), temporary.Sync(), temporary.Close())
}

func mode(original fs.FileInfo) fs.FileMode {
	if original == nil {
		return 0o644
	} else {
		return original.Mode() & (fs.ModePerm | fs.ModeSetuid | fs.ModeSetgid | fs.ModeSticky)
	}
}
