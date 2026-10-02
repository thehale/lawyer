// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"
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

func git(t *testing.T, date string, args ...string) {
	command := exec.Command("git", append([]string{"-c", "user.name=Test", "-c", "user.email=test@example.com", "-c", "commit.gpgsign=false"}, args...)...)
	command.Env = append(os.Environ(), "GIT_AUTHOR_DATE="+date, "GIT_COMMITTER_DATE="+date)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func commit(t *testing.T, year string, files map[string]string) {
	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git(t, year+"-06-01T12:00:00Z", "add", ".")
	git(t, year+"-06-01T12:00:00Z", "commit", "--quiet", "--message", year)
}

func createRepository(t *testing.T) {
	createProject(t, nil)
	_ = os.Remove("LICENSE")
	git(t, "", "init", "--quiet")
}

var this = strconv.Itoa(time.Now().Year())

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

var headed = text(`
# Copyright (c) 2026 Joseph Hale
# SPDX-License-Identifier: MPL-2.0
`)

func TestUsage(t *testing.T) {
	for _, args := range [][]string{{}, {"check"}, {"check", "--license", "MIT"}, {"check", "--bogus"}, {"check", "--license"}, append([]string{"check", "--fix=never"}, mine...), append([]string{"check", "--license", "MIT"}, mine...)} {
		if got := invoke(t, "", args...); got.code != 2 || !strings.Contains(got.stderr, "Usage:") {
			t.Errorf("%q: exit %d, stderr %q", args, got.code, got.stderr)
		}
	}
}

func TestWalkReportsEveryFailure(t *testing.T) {
	createProject(t, map[string]string{"a.sh": "echo\n", "b.yml": "x: 1\n", "c.sh": headed, "vendor/d.sh": "echo\n", "e.json": "{}\n", "COPYING.md": "text\n"})
	git(t, "", "init", "--quiet")

	got := invoke(t, "", append([]string{"check", "--exclude", "**/vendor/**"}, mine...)...)
	expected := text(`
a.sh: missing header
b.yml: missing header
`)

	if got.code != 1 || got.stderr != expected {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestStdinChecksOnlyTheListed(t *testing.T) {
	createProject(t, map[string]string{"a.sh": "echo\n", "b.sh": "echo\n"})
	git(t, "", "init", "--quiet")

	got := invoke(t, "b.sh\n", append([]string{"check", "-"}, mine...)...)

	if got.code != 1 || got.stderr != "b.sh: missing header\n" {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestAPathGivenTwiceIsCheckedOnce(t *testing.T) {
	createProject(t, map[string]string{"src/a.sh": "echo\n", "src/b.sh": "echo\n"})
	git(t, "", "init", "--quiet")

	got := invoke(t, "", append(append([]string{"check"}, mine...), "src/a.sh", "src", "./src/a.sh")...)
	expected := text(`
src/a.sh: missing header
src/b.sh: missing header
`)

	if got.code != 1 || got.stderr != expected {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestEmptyStdinPasses(t *testing.T) {
	createProject(t, map[string]string{"a.sh": "echo\n"})
	git(t, "", "init", "--quiet")

	if got := invoke(t, "", append([]string{"check", "-"}, mine...)...); got.code != 0 || got.stderr != "" {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestFixListsWhatChanged(t *testing.T) {
	createProject(t, map[string]string{"a.sh": "echo\n", "c.sh": headed})
	git(t, "", "init", "--quiet")

	fix := invoke(t, "", append([]string{"check", "--fix"}, mine...)...)
	again := invoke(t, "", append([]string{"check"}, mine...)...)

	if fix.code != 0 || fix.stdout != "a.sh: fixed\n" || again.code != 0 {
		t.Errorf("fix: exit %d, stdout %q; check after: exit %d, stderr %q", fix.code, fix.stdout, again.code, again.stderr)
	}
}

func TestMissingFileFails(t *testing.T) {
	createProject(t, nil)

	if got := invoke(t, "", append([]string{"check", "gone.sh"}, mine...)...); got.code != 1 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestYearsSpanTheEdits(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": "echo 1\n"})
	commit(t, "2025", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 1
`)})
	commit(t, "2026", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 2
`)})
	writeLicense(t)

	check := invoke(t, "", append([]string{"check"}, mine...)...)
	fix := invoke(t, "", append([]string{"check", "--fix"}, mine...)...)
	content, _ := os.ReadFile("a.sh")

	if check.stderr != "a.sh: years 2024, expected 2024-2026\n" || fix.code != 0 || !strings.HasPrefix(string(content), "# Copyright (c) 2024-2026 Joseph Hale\n") {
		t.Errorf("check: %q; fix: exit %d, %q; file: %q", check.stderr, fix.code, fix.stderr, content)
	}
}

func TestHeaderOnlyCommitsAreNotEdits(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": "echo 1\n"})
	commit(t, "2025", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 1
`)})
	writeLicense(t)

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.code != 0 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestUncommittedEditsCountAsThisYear(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 1
`)})
	project := map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 2
`), "new.sh": "echo\n"}
	for path, content := range project {
		_ = os.WriteFile(path, []byte(content), 0o644)
	}

	writeLicense(t)
	got := invoke(t, "", append([]string{"check"}, mine...)...)

	expected := fmt.Sprintf(text(`
a.sh: years 2024, expected 2024-%s
new.sh: missing header
`), this)

	if got.stderr != expected {
		t.Errorf("stderr %q", got.stderr)
	}
}

func TestShallowClonesOnlyCheckYearsArePast(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": "echo 1\n"})
	commit(t, "2026", map[string]string{"a.sh": text(`
# Copyright (c) 2020 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 2
`)})
	origin := must(os.Getwd())
	t.Chdir(t.TempDir())
	git(t, "", "clone", "--quiet", "--depth=1", "file://"+origin, ".")
	writeLicense(t)

	got := invoke(t, "", append([]string{"check"}, mine...)...)

	if got.code != 0 || got.stderr != "lawyer: this clone is shallow, so years are only checked for being in the past\n" {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestGivenYearsOverrideHistory(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 1
`)})

	writeLicense(t)
	got := invoke(t, "", append([]string{"check", "--copyright-year", "2020-2021"}, mine...)...)

	if got.stderr != "a.sh: years 2024, expected 2020-2021\n" {
		t.Errorf("stderr %q", got.stderr)
	}
}

func TestOutsideGitWarns(t *testing.T) {
	createProject(t, map[string]string{"a.sh": headed})

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.code != 0 || !strings.Contains(got.stderr, "not a git checkout") {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestFixWritesAMissingLicense(t *testing.T) {
	createProject(t, nil)
	_ = os.Remove("LICENSE")

	check := invoke(t, "", append([]string{"check"}, mine...)...)
	fix := invoke(t, "", append([]string{"check", "--fix"}, mine...)...)
	content, _ := os.ReadFile("LICENSE")
	top := "Mozilla Public License Version 2.0\n"

	if !strings.HasSuffix(check.stderr, "LICENSE: missing\n") || fix.stdout != "LICENSE: fixed\n" || !strings.HasPrefix(string(content), top) {
		t.Errorf("check %q; fix %q; LICENSE %q", check.stderr, fix.stdout, content[:min(80, len(content))])
	}
}

func TestSeveralLicensesTakeAFileEach(t *testing.T) {
	createProject(t, nil)
	_ = os.Remove("LICENSE")
	args := []string{"check", "--copyright-owner", "Joseph Hale", "--license", "MIT OR Apache-2.0"}

	fix := invoke(t, "", append(args, "--fix")...)
	again := invoke(t, "", args...)

	if fix.stdout != text(`
LICENSE-MIT: fixed
LICENSE-Apache-2.0: fixed
`) || again.code != 0 {
		t.Errorf("fix %q %q; check after: exit %d, %q", fix.stdout, fix.stderr, again.code, again.stderr)
	}
}

func TestSeveralLicensesReplaceThePlainOne(t *testing.T) {
	createProject(t, nil)
	args := []string{"check", "--copyright-owner", "Joseph Hale", "--license", "MIT OR Apache-2.0"}

	check := invoke(t, "", args...)
	fix := invoke(t, "", append(args, "--fix")...)
	again := invoke(t, "", args...)
	_, kept := os.Stat("LICENSE")

	if !strings.Contains(check.stderr, "LICENSE: not a file the --license expression calls for") || !strings.Contains(fix.stdout, "LICENSE: removed") || kept == nil || again.code != 0 {
		t.Errorf("check %q; fix %q; check after: exit %d, %q", check.stderr, fix.stdout, again.code, again.stderr)
	}
}

func TestAStrayLicenseFileIsRemoved(t *testing.T) {
	createProject(t, map[string]string{"LICENSE-MIT.md": "old\n"})

	fix := invoke(t, "", append([]string{"check", "--fix"}, mine...)...)
	_, kept := os.Stat("LICENSE-MIT.md")

	if fix.stdout != "LICENSE-MIT.md: removed\n" || kept == nil {
		t.Errorf("fix %q %q", fix.stdout, fix.stderr)
	}
}

func TestMarkdownLicensePasses(t *testing.T) {
	createProject(t, nil)
	_ = os.Rename("LICENSE", "LICENSE.md")

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.code != 0 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestFixLeavesAPassingLicenseAsWritten(t *testing.T) {
	createProject(t, nil)
	reflowed := strings.ReplaceAll(string(must(os.ReadFile("LICENSE"))), "\n\n", "\n\n  ")
	_ = os.WriteFile("LICENSE", []byte(reflowed), 0o644)

	got := invoke(t, "", append([]string{"check", "--fix"}, mine...)...)
	kept, _ := os.ReadFile("LICENSE")

	if got.code != 0 || got.stdout != "" || string(kept) != reflowed {
		t.Errorf("exit %d, stdout %q, LICENSE rewritten: %v", got.code, got.stdout, string(kept) != reflowed)
	}
}

func TestTheLicenseNamesTheOwnerAndLicense(t *testing.T) {
	createProject(t, map[string]string{"a.sh": "echo\n"})
	_ = os.Remove("LICENSE")
	git(t, "", "init", "--quiet")
	setup := invoke(t, "", "check", "--fix", "--copyright-owner", "Joseph Hale", "--license", "MIT")

	got := invoke(t, "", "check")

	if setup.code != 0 || got.code != 0 {
		t.Errorf("setup exit %d, stderr %q; check exit %d, stderr %q", setup.code, setup.stderr, got.code, got.stderr)
	}
}

func TestAChoiceOfLicensesIsWrittenAlphabetically(t *testing.T) {
	createProject(t, map[string]string{"a.sh": "echo\n"})
	_ = os.Remove("LICENSE")
	git(t, "", "init", "--quiet")
	fix := invoke(t, "", "check", "--fix", "--copyright-owner", "Joseph Hale", "--license", "MIT OR Apache-2.0")
	header, _ := os.ReadFile("a.sh")

	again := invoke(t, "", "check")

	if fix.code != 0 || !strings.Contains(string(header), "SPDX-License-Identifier: Apache-2.0 OR MIT") || again.code != 0 {
		t.Errorf("fix exit %d, a.sh %q; check without flags exit %d, stderr %q", fix.code, header, again.code, again.stderr)
	}
}

func TestALicenseWithoutAnOwnerLeavesTheOwnerRequired(t *testing.T) {
	createProject(t, nil)

	got := invoke(t, "", "check")

	if got.code != 2 || !strings.Contains(got.stderr, "an owner is required") || strings.Contains(got.stderr, "a license is required") {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestDuplicateLicenseFilesFail(t *testing.T) {
	createProject(t, map[string]string{"LICENSE.md": "copy\n"})

	if got := invoke(t, "", append([]string{"check"}, mine...)...); !strings.Contains(got.stderr, "LICENSE: one of several files for MPL-2.0\nLICENSE.md: one of several files for MPL-2.0\n") {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestUnknownLicenseIsAUsageError(t *testing.T) {
	createProject(t, nil)

	if got := invoke(t, "", "check", "--copyright-owner", "Joseph Hale", "--license", "CC-BY-4.0"); got.code != 2 || !strings.Contains(got.stderr, "unknown license CC-BY-4.0") {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestFixNeverFollowsALicenseSymlink(t *testing.T) {
	target := filepath.Join(t.TempDir(), "target")
	createProject(t, nil)
	_ = os.Remove("LICENSE")
	_ = os.WriteFile(target, []byte("keep\n"), 0o644)
	if err := os.Symlink(target, "LICENSE"); err != nil {
		t.Fatal(err)
	}

	got := invoke(t, "", append([]string{"check", "--fix"}, mine...)...)
	kept, _ := os.ReadFile(target)

	if string(kept) != "keep\n" || !strings.Contains(got.stderr, "LICENSE: not a regular file") {
		t.Errorf("target %q; exit %d, stderr %q", kept, got.code, got.stderr)
	}
}

func TestFixRefusesPathsOutsideTheWorkingDirectory(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "outside.sh")
	_ = os.WriteFile(outside, []byte("echo\n"), 0o644)
	createProject(t, nil)

	got := invoke(t, outside+"\n", append([]string{"check", "--fix", "-"}, mine...)...)
	kept, _ := os.ReadFile(outside)

	if got.code != 1 || string(kept) != "echo\n" || !strings.Contains(got.stderr, "outside the repository") {
		t.Errorf("exit %d, stderr %q, file %q", got.code, got.stderr, kept)
	}
}

func TestOnlyLicenseNamesSkipTheHeader(t *testing.T) {
	createProject(t, map[string]string{"LICENSE.sh": "echo\n", "COPYING-MIT.txt": "text\n"})

	if got := invoke(t, "", append([]string{"check", "LICENSE.sh", "COPYING-MIT.txt"}, mine...)...); !strings.Contains(got.stderr, "LICENSE.sh: missing header") {
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

func TestControlCharactersInNamesAreQuoted(t *testing.T) {
	createProject(t, map[string]string{"evil\x1b[2J.sh": "echo\n"})

	if got := invoke(t, "", append([]string{"check"}, mine...)...); !strings.Contains(got.stderr, `"evil\x1b[2J.sh": missing header`) {
		t.Errorf("stderr %q", got.stderr)
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

func TestCommentsThatMentionCopyrightAreEdits(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

# Copyright (c) notices are printed in one pass
`)})
	commit(t, "2025", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

# Copyright (c) notices are printed in two passes
`)})
	writeLicense(t)

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.stderr != "a.sh: years 2024, expected 2024-2025\n" {
		t.Errorf("stderr %q", got.stderr)
	}
}

func TestFutureCommitDatesAreIgnored(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 1
`)})
	commit(t, "2099", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 2
`)})
	writeLicense(t)

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.code != 0 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestAnOwnerSpanningLinesIsAUsageError(t *testing.T) {
	createProject(t, nil)

	if got := invoke(t, "", "check", "--copyright-owner", "Acme\nrm -rf /", "--license", "MIT"); got.code != 2 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestARefusedPathLeavesTheRestChecked(t *testing.T) {
	createProject(t, map[string]string{"bad.sh": "echo\n"})

	got := invoke(t, text(`
bad.sh
../outside.sh
`), append([]string{"check", "-"}, mine...)...)

	if got.code != 1 || !strings.Contains(got.stderr, "bad.sh: missing header") || !strings.Contains(got.stderr, "outside the repository") {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestMergeResolutionsAreEdits(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": text(`
# Copyright (c) 2024-2025 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo base
`)})
	git(t, "", "branch", "--quiet", "side")
	commit(t, "2025", map[string]string{"a.sh": text(`
# Copyright (c) 2024-2025 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo main
`)})
	git(t, "", "checkout", "--quiet", "side")
	commit(t, "2025", map[string]string{"a.sh": text(`
# Copyright (c) 2024-2025 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo side
`)})
	git(t, "", "checkout", "--quiet", "-")
	merge := exec.Command("git", "-c", "user.name=Test", "-c", "user.email=test@example.com", "merge", "--quiet", "side")
	_ = merge.Run()
	commit(t, "2026", map[string]string{"a.sh": text(`
# Copyright (c) 2024-2025 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo resolved
`)})
	writeLicense(t)

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.stderr != "a.sh: years 2024-2025, expected 2024-2026\n" {
		t.Errorf("stderr %q", got.stderr)
	}
}

func TestProseAfterAnIdentifierIsAnEdit(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

# SPDX-License-Identifier: MIT needs a note
`)})
	commit(t, "2025", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

# SPDX-License-Identifier: MIT needs a clearer note
`)})
	writeLicense(t)

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.stderr != "a.sh: years 2024, expected 2024-2025\n" {
		t.Errorf("stderr %q", got.stderr)
	}
}

func TestSourcesMarkedBinaryKeepTheirHistory(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{".gitattributes": "*.sh binary\n", "a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 1
`)})
	commit(t, "2025", map[string]string{"a.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 2
`)})
	writeLicense(t)

	if got := invoke(t, "", append([]string{"check"}, mine...)...); !strings.Contains(got.stderr, "a.sh: years 2024, expected 2024-2025") {
		t.Errorf("stderr %q", got.stderr)
	}
}

func TestNamesLikePathspecMagicKeepTheirHistory(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{":(exclude)odd.sh": text(`
# Copyright (c) 2024 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

echo 1
`)})
	writeLicense(t)

	if got := invoke(t, "", append([]string{"check"}, mine...)...); got.code != 0 {
		t.Errorf("exit %d, stderr %q", got.code, got.stderr)
	}
}

func TestLateNullBytesStillMarkABinary(t *testing.T) {
	createProject(t, map[string]string{"asset.sh": strings.Repeat("a", 9000) + "\x00payload"})

	if got := invoke(t, "", append([]string{"check", "--fix", "asset.sh"}, mine...)...); got.stdout != "" {
		t.Errorf("stdout %q", got.stdout)
	}
}

func TestRefusedPathsAreQuoted(t *testing.T) {
	createProject(t, nil)

	if got := invoke(t, "../\x1b[2Jx.sh\n", append([]string{"check", "-"}, mine...)...); !strings.Contains(got.stderr, `"../\x1b[2Jx.sh": outside`) {
		t.Errorf("stderr %q", got.stderr)
	}
}

func TestRewritingANoticeIsNotAnEdit(t *testing.T) {
	createRepository(t)
	commit(t, "2024", map[string]string{"a.sh": text(`
# Copyright (c) 2024 - 2024 Joseph Hale
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.

echo
`)})
	writeLicense(t)

	fix := invoke(t, "", append([]string{"check", "--fix"}, mine...)...)
	again := invoke(t, "", append([]string{"check"}, mine...)...)
	content, _ := os.ReadFile("a.sh")

	if fix.code != 0 || again.code != 0 || !strings.HasPrefix(string(content), "# Copyright (c) 2024 Joseph Hale\n# SPDX-License-Identifier: MPL-2.0\n\necho") {
		t.Errorf("fix: %q; check after: %q; file: %q", fix.stderr, again.stderr, content)
	}
}

func TestLanguagesListsEachWithItsFiles(t *testing.T) {
	got := invoke(t, "", "languages")
	table := regexp.MustCompile(` {2,}`).ReplaceAllString(got.stdout, "  ")

	if got.code != 0 || !strings.Contains(table, "Bash  .bash .sh #!bash #!sh\n") || !strings.Contains(table, "Ruby  Gemfile .rb .gemspec .rake .ru #!ruby\n") || !strings.Contains(table, "JSON  .json (no header)\n") {
		t.Errorf("exit %d, stdout %q", got.code, got.stdout)
	}
}

func TestLicensesListsTheAcceptedIDs(t *testing.T) {
	got := invoke(t, "", "licenses")
	ids := strings.Fields(got.stdout)

	if got.code != 0 || !slices.Contains(ids, "MPL-2.0") || !slices.Contains(ids, "GPL-3.0-or-later") || slices.Contains(ids, "CC-BY-4.0") {
		t.Errorf("exit %d, stdout %q", got.code, got.stdout)
	}
}
