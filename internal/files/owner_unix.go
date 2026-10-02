// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

//go:build unix

package files

import (
	"io/fs"
	"os"
	"syscall"
)

func own(temporary *os.File, original fs.FileInfo) error {
	was, wasKnown := sys(original)
	now, nowKnown := sys(info(temporary))

	if wasKnown && nowKnown && (was.Uid != now.Uid || was.Gid != now.Gid) {
		return temporary.Chown(int(was.Uid), int(was.Gid))
	} else {
		return nil
	}
}

func sys(info fs.FileInfo) (*syscall.Stat_t, bool) {
	if info == nil {
		return nil, false
	} else {
		stat, ok := info.Sys().(*syscall.Stat_t)
		return stat, ok
	}
}

func info(file *os.File) fs.FileInfo {
	info, err := file.Stat()
	if err == nil {
		return info
	} else {
		return nil
	}
}
