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

func ByID(id ID) (License, bool) {
	text, textErr := texts.ReadFile("texts/" + string(id) + ".txt")
	source, templateErr := texts.ReadFile("texts/" + string(id) + ".template.txt")
	return License{ID: id, Text: string(text), template: templateOf(string(source))}, textErr == nil && templateErr == nil
}

func IDs() []ID {
	templates, _ := fs.Glob(texts, "texts/*.template.txt")
	var ids []ID
	for _, file := range templates {
		ids = append(ids, ID(strings.TrimSuffix(path.Base(file), ".template.txt")))
	}
	return ids
}
