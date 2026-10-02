#!/usr/bin/env bats
# Copyright (c) 2026 Joseph Hale
# SPDX-License-Identifier: MPL-2.0

bats_require_minimum_version 1.5.0

setup() {
	REPO="$(cd "$BATS_TEST_DIRNAME/../.." && pwd)"
}

@test "takes no arguments, since the licenses are SPDX's OSI-approved set" {
	run --separate-stderr "$REPO/bin/update-licenses" MIT

	[ "$status" -eq 2 ]
	[[ "$stderr" == "Usage: bin/update-licenses" ]]
}
