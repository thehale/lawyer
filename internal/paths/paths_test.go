// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package paths

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

func TestWalkSeesWhatGitSees(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	runGit(t, "init", "--quiet")
	write(t, ".gitignore", "/build/\n")
	write(t, "tracked.sh", "")
	write(t, "untracked.sh", "")
	write(t, "build/output.sh", "")
	write(t, "vendor/lib.sh", "")
	runGit(t, "add", ".gitignore", "tracked.sh")

	files, err := FilesUnder(".")
	slices.Sort(files)

	if err != nil || !slices.Equal(files, []Path{".gitignore", "tracked.sh", "untracked.sh", "vendor/lib.sh"}) {
		t.Errorf("got %q, %v", files, err)
	}
}

func TestWalkOutsideGit(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(must(os.Getwd())))
	write(t, "a.sh", "")
	write(t, "sub/b.sh", "")

	files, err := FilesUnder(".")

	if err != nil || !slices.Equal(files, []Path{"a.sh", Path("sub").Child("b.sh")}) {
		t.Errorf("got %q, %v", files, err)
	}
}

func runGit(t *testing.T, args ...string) {
	if output, err := exec.Command("git", args...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func write(t *testing.T, path, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func must(value string, err error) string {
	if err != nil {
		panic(err)
	}
	return value
}
