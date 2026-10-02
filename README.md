<!--
Copyright (c) 2026 Joseph Hale
SPDX-License-Identifier: MPL-2.0
-->

<div align="center">

# lawyer

A linter for your copyright headers (with autofix!)

<!-- BADGES -->
[![License: MPL-2.0](https://badgen.net/github/license/thehale/lawyer)](https://github.com/thehale/lawyer/blob/main/LICENSE)
[![Sponsor thehale on GitHub](https://badgen.net/badge/icon/Sponsor/pink?icon=github&label)](https://github.com/sponsors/thehale)
[![Joseph Hale's software engineering blog](https://jhale.dev/badges/website.svg)](https://jhale.dev)
[![Follow Joseph Hale on LinkedIn](https://jhale.dev/badges/follow.svg)](https://www.linkedin.com/comm/mynetwork/discovery-see-all?usecase=PEOPLE_FOLLOWS&followMember=thehale)
</div>

## Quickstart

```bash
lawyer check --copyright-owner "Your Name" --license MPL-2.0
```

Every file gets a header in its own comment syntax:

```bash
# Copyright (c) 2024-2026 Your Name
# SPDX-License-Identifier: MPL-2.0
```

- **Files.** With no paths, lawyer checks what git lists: tracked files and
  untracked files that aren't ignored. `-` reads the paths from stdin instead.
  `--exclude '**/glob/**'` checks nothing on the paths it matches.
- **Existing headers.** The owner's copyright line is found in any common
  shape, such as `Copyright 2022, 2024 Your Name` or `© 2022 Your Name`, and
  rewritten to the one above. A license's standard notice, such as MPL-2.0's
  "This Source Code Form is subject to…", becomes the SPDX line. Other owners'
  copyright lines stay.
- **Years.** A header spans the years the file was edited, from git history.
  Commits that only touch the header don't count, and uncommitted changes count
  as this year. A shallow clone or a directory outside git is checked less
  strictly, with a warning. `--copyright-year 2024-2026` overrides the history.
- **LICENSE.** It must be the license's text, with the owner's copyright line
  wherever the license has a place for one, as MIT and BSD do.
  An expression with several licenses, such as `MIT OR Apache-2.0`, takes a
  `LICENSE-<id>` file for each.

`lawyer languages` lists the languages lawyer puts headers in, and
`lawyer licenses` lists the SPDX ids `--license` accepts: every license that
is both on SPDX's list and approved by the OSI.

## Development

```bash
bin/setup  # Install the tools
bin/build  # Build lawyer into build/lawyer
bin/ci     # Run the checks
bin/ci --fix  # Fix what can be fixed automatically
```

Each language has its own file in [internal/languages](internal/languages).
License texts come from SPDX into
[internal/spdx/texts](internal/spdx/texts) through `bin/update-licenses`.
The [add-a-language](.agents/skills/add-a-language/SKILL.md) skill walks
through adding a language.

## License

Copyright (c) 2026 Joseph Hale, All Rights Reserved

Provided under the terms of the [Mozilla Public License, version 2.0](./LICENSE)

<details>

<summary><b>What does the MPL-2.0 license allow/require?</b></summary>

### TL;DR

You can use files from this project in both open source and proprietary
applications, provided you include the above attribution. However, if
you modify any code in this project, or copy blocks of it into your own
code, you must publicly share the resulting files (note, not your whole
program) under the MPL-2.0. The best way to do this is via a Pull
Request back into this project.

If you have any other questions, you may also find Mozilla's [official
FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/) for the MPL-2.0 license
insightful.

If you dislike this license, you can contact me about negotiating a paid
contract with different terms.

**Disclaimer:** This TL;DR is just a summary. All legal questions
regarding usage of this project must be handled according to the
official terms specified in the `LICENSE` file.

### Why the MPL-2.0 license?

I believe that an open-source software license should ensure that code
can be used everywhere.

Strict copyleft licenses, like the GPL family of licenses, fail to
fulfill that vision because they only permit code to be used in other
GPL-licensed projects. Permissive licenses, like the MIT and Apache
licenses, allow code to be used everywhere but fail to prevent
proprietary or GPL-licensed projects from limiting access to any
improvements they make.

In contrast, the MPL-2.0 license allows code to be used in any software
project, while ensuring that any improvements remain available for
everyone.

</details>
