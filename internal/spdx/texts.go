// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package spdx

import (
	"embed"
	"io/fs"
	"path"
	"strings"
)

//go:embed texts/*.txt
var texts embed.FS

func IDs() []ID {
	templates, _ := fs.Glob(texts, "texts/*.template.txt")
	var ids []ID
	for _, file := range templates {
		ids = append(ids, ID(strings.TrimSuffix(path.Base(file), ".template.txt")))
	}
	return ids
}
