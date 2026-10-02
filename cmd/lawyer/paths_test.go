// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/thehale/lawyer"
)

func TestPathsReadStdinLines(t *testing.T) {
	files, err := paths([]string{"first.sh", "-"}, strings.NewReader("a.sh\r\n\nb.sh\n"))

	if err != nil || !slices.Equal(files, []lawyer.Path{"first.sh", "a.sh", "b.sh"}) {
		t.Errorf("got %q, %v", files, err)
	}
}
