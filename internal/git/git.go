// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package git

import (
	"os"
	"os/exec"
	"strings"
)

// Run runs git in dir, taking every pathspec literally, and returns what it
// printed.
func Run(dir string, args ...string) (string, error) {
	command := exec.Command("git", append([]string{"-C", dir}, args...)...)
	command.Env = append(os.Environ(), "GIT_LITERAL_PATHSPECS=1")
	output, err := command.Output()
	return string(output), err
}

// IsInCheckout reports whether dir is inside a git checkout.
func IsInCheckout(dir string) bool {
	output, err := Run(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(output) == "true"
}
