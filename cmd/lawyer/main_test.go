// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binary string

func TestMain(m *testing.M) {
	directory, err := os.MkdirTemp("", "lawyer")
	binary = filepath.Join(directory, "lawyer")
	if err == nil {
		err = exec.Command("go", "build", "-o", binary, ".").Run()
	}
	if err != nil {
		panic(err)
	}
	code := m.Run()
	_ = os.RemoveAll(directory)
	os.Exit(code)
}

type result struct {
	code           int
	stdout, stderr string
}

func invoke(t *testing.T, stdin string, args ...string) result {
	command := exec.Command(binary, args...)
	command.Stdin = strings.NewReader(stdin)
	var stdout, stderr strings.Builder
	command.Stdout, command.Stderr = &stdout, &stderr
	err := command.Run()
	if _, exited := err.(*exec.ExitError); err != nil && !exited {
		t.Fatal(err)
	}
	return result{code: command.ProcessState.ExitCode(), stdout: stdout.String(), stderr: stderr.String()}
}

func createProject(t *testing.T, files map[string]string) {
	t.Chdir(t.TempDir())
	t.Setenv("GIT_CEILING_DIRECTORIES", filepath.Dir(must(os.Getwd())))
	writeLicense(t)
	for path, content := range files {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

var mpl = must(os.ReadFile(filepath.Join("..", "..", "LICENSE")))

func writeLicense(t *testing.T) {
	if err := os.WriteFile("LICENSE", mpl, 0o644); err != nil {
		t.Fatal(err)
	}
}

func text(example string) string {
	return strings.TrimPrefix(example, "\n")
}

func must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

var mine = []string{"--copyright-owner", "Joseph Hale", "--license", "MPL-2.0"}

func TestUsage(t *testing.T) {
	for _, args := range [][]string{{}, {"check"}, {"check", "--license", "MIT"}, {"check", "--bogus"}, {"check", "--license"}, append([]string{"check", "--fix=never"}, mine...), append([]string{"check", "--license", "MIT"}, mine...)} {
		if got := invoke(t, "", args...); got.code != 2 || !strings.Contains(got.stderr, "Usage:") {
			t.Errorf("%q: exit %d, stderr %q", args, got.code, got.stderr)
		}
	}
}

func TestMissingFileFails(t *testing.T) {
	createProject(t, nil)

	if got := invoke(t, "", append([]string{"check", "gone.sh"}, mine...)...); got.code != 1 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestShebangAfterAByteOrderMark(t *testing.T) {
	createProject(t, map[string]string{"tool": "\ufeff" + text(`
#!/usr/bin/env python3
print()
`)})

	if got := invoke(t, "", append([]string{"check", "tool"}, mine...)...); !strings.Contains(got.stderr, "tool: missing header") {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestFixKeepsTheFileMode(t *testing.T) {
	createProject(t, map[string]string{"run.sh": "echo\n"})
	_ = os.Chmod("run.sh", 0o750)

	invoke(t, "", append([]string{"check", "--fix", "run.sh"}, mine...)...)
	info, _ := os.Stat("run.sh")

	if info.Mode().Perm() != 0o750 {
		t.Errorf("mode %v", info.Mode())
	}
}

func TestAFailedWriteIsNotReportedAsFixed(t *testing.T) {
	createProject(t, map[string]string{"locked/a.sh": "echo\n"})
	_ = os.Chmod("locked", 0o555)
	t.Cleanup(func() { _ = os.Chmod("locked", 0o755) })
	if os.WriteFile("locked/probe", nil, 0o644) == nil {
		t.Skip("the directory is still writable, as it is for root")
	}

	got := invoke(t, "", append([]string{"check", "--fix", "locked/a.sh"}, mine...)...)
	kept, _ := os.ReadFile("locked/a.sh")

	if got.code != 1 || got.stdout != "" || string(kept) != "echo\n" {
		t.Errorf("exit %d, stdout %q, file %q", got.code, got.stdout, kept)
	}
}

func TestAnOwnerSpanningLinesIsAUsageError(t *testing.T) {
	createProject(t, nil)

	if got := invoke(t, "", "check", "--copyright-owner", "Acme\nrm -rf /", "--license", "MIT"); got.code != 2 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestLateNullBytesStillMarkABinary(t *testing.T) {
	createProject(t, map[string]string{"asset.sh": strings.Repeat("a", 9000) + "\x00payload"})

	if got := invoke(t, "", append([]string{"check", "--fix", "asset.sh"}, mine...)...); got.stdout != "" {
		t.Errorf("stdout %q", got.stdout)
	}
}
