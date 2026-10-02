---
name: add-a-language
description: >-
  Add a language lawyer can put headers in. Use when lawyer skips a file type
  it should check.
license: MPL-2.0
---
<!--
Copyright (c) 2026 Joseph Hale
SPDX-License-Identifier: MPL-2.0
-->

# Add a language

Each language is one file in `internal/languages`, and nothing else lives
there. The types it's built from (`Language`, `Comment`, `Prolog`) live in
`internal/syntax`.

1. Copy the closest existing language to `internal/languages/<name>.go` and
   register it:

   ```go
   func init() {
       syntax.Register(syntax.Language{
           Name:       "Ruby",
           Filenames:  []string{"Gemfile"},
           Extensions: []string{"rb", "gemspec", "rake", "ru"},
           Shebangs:   []string{"ruby"},
           Comment:    syntax.Hashes,
           Prolog:     []syntax.Prolog{syntax.Shebang, syntax.Rack, syntax.Encoding},
       })
   }
   ```

   - `Filenames` are whole file names, such as `Dockerfile`, for files with no
     telling extension. `Extensions` go without the dot and may have several
     parts, as `cmake.in` does. Both match in any case. `Shebangs` names the
     interpreter, as it appears after `#!/usr/bin/env` or the last `/`. List
     the same names in the same order in both where they overlap.
   - `Comment` picks the style of a new header from `internal/syntax/comment.go`.
     Choose the form the language's own projects use for license headers: a
     block comment such as `Stars` where that is the convention, and line
     comments such as `Slashes` or `Hashes` where it isn't. Add a style there
     only when none fits.
   - `Prolog` lists the lines that must stay above the header, in order, from
     `internal/syntax/prolog.go`, such as `Shebang`, `Encoding` or
     `FrontMatter`. Add one there when the language has a line that must come
     first, such as CSS's `@charset`.
   - A language with no comment syntax registers without a `Comment`, and a
     comment above `init` says why, as `json.go` does. lawyer then knows the
     file and leaves it alone.

2. Add an example to `internal/languages/testdata/<name>/`, named for the
   file it stands for with `.txt` added, such as `main.c.txt`. It holds the
   file before `--fix`, a line of 80 dashes, and the file after. Add one for
   each prolog line, and one without the dashes for a header that files in
   the wild already carry in another comment style, which must pass as it is.

3. Run `bin/ci`. lawyer checks this repository with its own build, so a new
   language can flag files here that never had a header. `bin/ci --fix` adds
   them.
