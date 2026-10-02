// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package files

import (
	"os"
	"testing"

	"github.com/thehale/lawyer/internal/paths"
)

func TestReplaceLeavesAFileEditedSinceItWasRead(t *testing.T) {
	path := temporaryDirectory(t).Child("a.sh")
	save(t, path, "saved meanwhile\n")

	err := Replace(path, "read earlier\n", "fixed\n")

	if err == nil || contents(path) != "saved meanwhile\n" {
		t.Errorf("error %v, file %q", err, contents(path))
	}
}

func TestReplaceNeverRecreatesADeletedFile(t *testing.T) {
	path := temporaryDirectory(t).Child("empty.sh")

	err := Replace(path, "", "fixed\n")

	if err == nil || exists(path) {
		t.Errorf("error %v, recreated %v", err, exists(path))
	}
}

func temporaryDirectory(t *testing.T) paths.Path {
	return paths.Path(t.TempDir())
}

func save(t *testing.T, path paths.Path, content string) {
	if err := os.WriteFile(string(path), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func contents(path paths.Path) string {
	content, _ := os.ReadFile(string(path))
	return string(content)
}

func exists(path paths.Path) bool {
	_, err := os.Lstat(string(path))
	return err == nil
}
