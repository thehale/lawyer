// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package lawyer

import (
	"bytes"
	"os"

	"github.com/thehale/lawyer/internal/files"
	"github.com/thehale/lawyer/internal/header"
	_ "github.com/thehale/lawyer/internal/languages"
	"github.com/thehale/lawyer/internal/syntax"
)

// A Licensable is a file whose header must carry a Declaration, such as a
// source file.
type Licensable struct {
	repository Repository
	path       Path
	content    string
	header     Header
}

// A Header is the copyright line and license identifier or notice at the top
// of a Licensable.
type Header = header.Header

func licensableAt(repository Repository, path Path) (Licensable, bool, error) {
	content, language, err := readSource(path)

	if err == nil && language.HasHeader() {
		return Licensable{repository: repository, path: path, content: content, header: header.Read(language, content)}, true, nil
	} else {
		return Licensable{}, false, err
	}
}

// Path is where the licensable is.
func (l Licensable) Path() Path {
	return l.path
}

// Header is the licensable's header as it stands.
func (l Licensable) Header() Header {
	return l.header
}

// Violations are how the licensable's header falls short of declaration.
func (l Licensable) Violations(declaration Declaration) []Violation {
	return violationsAt(l.path, l.header.Violations(declaration.headerExpectation(l))...)
}

// ReplaceHeader writes replacement in place of the licensable's header.
func (l Licensable) ReplaceHeader(replacement Header) (Changes, error) {
	return fixOf(l.path, files.Replace(l.path, l.content, l.header.ContentWith(replacement)))
}

// Fix replaces the licensable's header with the one declaration calls for.
func (l Licensable) Fix(declaration Declaration) (Changes, error) {
	return l.ReplaceHeader(declaration.HeaderFor(l))
}

func readSource(path Path) (string, syntax.Language, error) {
	info, err := os.Lstat(string(path))

	if err == nil && info.Mode().IsRegular() {
		return readText(path)
	} else {
		return "", syntax.Language{}, err
	}
}

func readText(path Path) (string, syntax.Language, error) {
	content, err := os.ReadFile(string(path))
	language, _ := syntax.LanguageOf(path.Base(), string(content))

	if err == nil && !isBinary(content) {
		return string(content), language, nil
	} else {
		return "", syntax.Language{}, err
	}
}

func isBinary(content []byte) bool {
	return bytes.IndexByte(content, 0) >= 0
}
