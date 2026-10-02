---
name: add-a-license
description: >-
  Add a license lawyer can check. Use when lawyer refuses a --license id as
  unknown.
license: MPL-2.0
---
<!--
Copyright (c) 2026 Joseph Hale
SPDX-License-Identifier: MPL-2.0
-->

# Add a license

License texts come verbatim from SPDX, two or three files per license in
`internal/spdx/texts/`:

- `<id>.txt` is the canonical text that `--fix` writes.
- `<id>.template.txt` marks what may vary, and the check matches against it.
- `<id>.notice.txt` is the license's standard header notice, when SPDX has
  one. `--fix` replaces a notice it finds in a file's header with the SPDX
  line.

lawyer carries every license that SPDX marks OSI-approved and not deprecated,
the overlap of SPDX's list with <https://opensource.org/licenses>.
`bin/update-licenses` reads that set from SPDX's `licenses.json` and replaces
`internal/spdx/texts/` with it, so there is no per-license step:

1. Change `VERSION` in `bin/update-licenses` to the SPDX release that adds
   the license.
2. Run `bin/update-licenses`. It downloads the files for every license in the
   set, and refuses anything that isn't shaped like an id.
3. Run `bin/ci`. `TestRenderedLicensesPass` renders every license with an
   owner and checks it against its own template. A new license that fails
   means its template uses markup the matcher in `internal/spdx/match.go`
   doesn't handle yet.

A license outside that set, such as a non-OSI one, would need the set itself
to change, which is a question for the operator rather than a data update.

Expressions with `WITH` exceptions and `LicenseRef-` licenses aren't supported
yet, so an id of either kind needs code, not just data.
