// Copyright (c) 2026 Joseph Hale
// SPDX-License-Identifier: MPL-2.0

package header

import (
	"fmt"
	"slices"

	"github.com/thehale/lawyer/internal/copyright"
	"github.com/thehale/lawyer/internal/detect"
	"github.com/thehale/lawyer/internal/spdx"
	"github.com/thehale/lawyer/internal/syntax"
)

// A Header is the copyright line and license identifier or notice at the top
// of a file, and the comment lines they sit in.
type Header struct {
	language syntax.Language
	doc      document
	start    int
	comment  []syntax.Line
	claims   detect.Claims
	owner    string
}

// An Expectation is what a file's header must say: its copyright line and its
// license.
type Expectation struct {
	Copyright copyright.Expectation
	License   spdx.Expression
}

const identityKey = "SPDX-License-Identifier: "

// Read finds the header in content.
func Read(language syntax.Language, content string) Header {
	doc := documentOf(content)
	start := language.HeaderStart(doc.lines)
	comment := language.Comment.Read(doc.lines[start:])
	removable := comment[:min(len(comment), len(doc.lines)-1-start)]
	return Header{language: language, doc: doc, start: start, comment: removable, claims: detect.ClaimsIn(texts(comment))}
}

// IsCanonical reports whether the header is a copyright line and an identifier
// line alone, written in its language's standard form.
func (h Header) IsCanonical() bool {
	return h.copyrightOf("") != nil && h.identity() != nil && h.notice() == nil && h.region("").isExactly(h.Lines())
}

// Canonical is the header expectation calls for, keeping what this header's
// copyright line says where expectation allows it.
func (h Header) Canonical(expectation Expectation) Header {
	owner := expectation.Copyright.Owner
	line := statement(h.copyrightOf(owner)).Canonical(expectation.Copyright).String()
	claims := detect.Claims{detect.Copyright{Line: line}, detect.Identifier{License: expectation.License}}
	return Header{language: h.language, claims: claims, owner: owner}
}

// Lines are the header written in its language's standard form.
func (h Header) Lines() []string {
	return h.language.Comment.Render([]string{h.copyrightOf("").Line, identityKey + string(h.identity().License)})
}

// ContentWith is the file the header was read from, with replacement in its
// place.
func (h Header) ContentWith(replacement Header) string {
	lines := replacement.Lines()
	region := h.region(replacement.owner)

	if region.isExactly(lines) {
		return h.doc.String()
	} else {
		return h.doc.withReplaced(region.start, region.end(), region.remainingLines()).withInserted(region.start, lines).String()
	}
}

// Violations are how the header falls short of expectation.
func (h Header) Violations(expectation Expectation) []string {
	owner := expectation.Copyright.Owner
	violations := slices.Concat(h.copyrightViolations(expectation), h.identityViolations(expectation), h.noticeViolations(expectation))

	switch {
	case h.ContentWith(h.Canonical(expectation)) == h.doc.String():
		return nil
	case h.copyrightOf(owner) == nil && h.identity() == nil && h.notice() == nil:
		return []string{"missing header"}
	case violations == nil:
		return []string{"header isn't in the standard form"}
	default:
		return violations
	}
}

func (h Header) copyrightOf(owner string) *detect.Copyright {
	return first(h.claims, namesOwner(owner))
}

func (h Header) identity() *detect.Identifier {
	return first(h.claims, anyClaim[detect.Identifier])
}

func (h Header) notice() *detect.Notice {
	return first(h.claims, anyClaim[detect.Notice])
}

func (h Header) region(owner string) region {
	return region{start: h.start, lines: regionLines(h.doc.lines[h.start:], h.comment, h.claims, namesOwner(owner))}
}

func (h Header) copyrightViolations(expectation Expectation) []string {
	claim := h.copyrightOf(expectation.Copyright.Owner)
	years := expectation.Copyright.Violations(statement(claim))
	canonical := statement(claim).Canonical(expectation.Copyright).String()

	switch {
	case claim == nil:
		return []string{fmt.Sprintf("missing copyright line for %s", expectation.Copyright.Owner)}
	case years != nil:
		return years
	case claim.Line != canonical:
		return []string{fmt.Sprintf("copyright line should read %q", canonical)}
	default:
		return nil
	}
}

func (h Header) identityViolations(expectation Expectation) []string {
	identity := h.identity()

	switch {
	case identity == nil && h.notice() != nil:
		return nil
	case identity == nil:
		return []string{"missing SPDX-License-Identifier"}
	case identity.License != expectation.License:
		return []string{fmt.Sprintf("expected %s, found %s", expectation.License, identity.License)}
	default:
		return nil
	}
}

func (h Header) noticeViolations(expectation Expectation) []string {
	notice := h.notice()

	switch {
	case notice == nil:
		return nil
	case notice.License == expectation.License:
		return []string{fmt.Sprintf("%s notice in place of an SPDX-License-Identifier line", notice.License)}
	default:
		return []string{fmt.Sprintf("expected %s, found the %s notice", expectation.License, notice.License)}
	}
}

func statement(claim *detect.Copyright) copyright.Statement {
	if claim == nil {
		return ""
	} else {
		return claim.Statement
	}
}

func namesOwner(owner string) func(detect.Copyright) bool {
	return func(claim detect.Copyright) bool { return claim.Statement.HasOwner(owner) }
}

func texts(comment []syntax.Line) []string {
	texts := make([]string, len(comment))
	for index, line := range comment {
		texts[index] = line.Text
	}
	return texts
}

func regionLines(raw []string, comment []syntax.Line, claims detect.Claims, namesOwner func(detect.Copyright) bool) []regionLine {
	foreign := func(claim detect.Claim) bool {
		copyright, ok := claim.(detect.Copyright)
		return ok && !namesOwner(copyright)
	}
	mine := slices.DeleteFunc(slices.Clone(claims), foreign)
	others := slices.DeleteFunc(slices.Clone(claims), func(claim detect.Claim) bool { return !foreign(claim) })
	lines := make([]regionLine, len(comment))
	for index, line := range comment {
		lines[index] = regionLine{raw: raw[index], syntax: line.Syntax, ours: mine.HasLine(index) && !others.HasLine(index)}
	}
	return lines
}

func anyClaim[T detect.Claim](T) bool {
	return true
}

func first[T detect.Claim](claims detect.Claims, matches func(T) bool) *T {
	index := slices.IndexFunc(claims, func(claim detect.Claim) bool {
		typed, ok := claim.(T)
		return ok && matches(typed)
	})

	if index >= 0 {
		typed := claims[index].(T)
		return &typed
	} else {
		return nil
	}
}
