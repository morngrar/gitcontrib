// Copyright 2023 gitcontrib Authors
// SPDX-License-Identifier: Apache-2.0

// Package example provides the Bonzai command branch of the same name.
package gitcontrib

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/conf"
	"github.com/rwxrob/help"
	"github.com/rwxrob/vars"
)

// Most Cmds that make use of the conf and vars branches will want to
// call SoftInit in order to create the persistence layers or whatever
// else is needed to initialize their use. This cannot be done
// automatically from these imported modules because Cmd authors may
// with to change the default values before calling SoftInit and
// committing them.

func init() {
	Z.Conf.SoftInit()
	Z.Vars.SoftInit()
}

// Cmd provides a Bonzai branch command that can be composed into Bonzai
// trees or used as a standalone with light wrapper (see cmd/).
var Cmd = &Z.Cmd{

	Name:      `gitcontrib`,
	Summary:   `a command tree branch for analysing git author contributions`,
	Version:   `v0.0.1`,
	Copyright: `Copyright 2023 Svein-Kåre Bjørnsen`,
	License:   `Apache-2.0`,
	Source:    `git@github.com:PLACEHOLDER/gitcontrib.git`,
	Issues:    `github.com/PLACEHOLDER/gitcontrib/issues`,

	// Composite commands, local and external, all have their own names
	// that are added to the command tree depending on where they are
	// composed.

	Commands: []*Z.Cmd{

		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd, conf.Cmd, vars.Cmd,

		// local commands (in this module)
		gitStatusCmd, GitAuthorCommitsCmd, GitMapAuthorChangesCmd,
		ContributionSummaryCmd,
	},

	// Add custom BonzaiMark template extensions (or overwrite existing ones).
	Dynamic: template.FuncMap{
		"uname": func(_ *Z.Cmd) string { return Z.Out("uname", "-a") },
		"dir":   func() string { return Z.Out("dir") },
		"pwd":   func() string { return Z.Out("pwd") },
	},

	// WARNING: The Description will be dedented using the exact
	// runes from the beginning of the first line to the first
	// non-whitespace character. This means that mixing spaces and tabs in
	// your indentation will created unwanted truncation. Make sure each
	// line has the same indentation exactly.

	Description: `
		The {{aka}} command is a well-documented example to get you started.
		You can start the description here and wrap it to look nice and it
		will just work.  Descriptions are written in BonzaiMark,
		a simplified combination of CommonMark, "go doc", and text/template
		that uses the Cmd itself as a data source and has a rich set of
		builtin template functions ({{pre "pre"}}, {{pre "exename"}},
		{{pre "indent"}}, etc.). There are four block types and four span types in
		BonzaiMark:

		Spans

		    Plain
		    *Italic*
		    **Bold**
		    ***BoldItalic***
		    <Under> (brackets remain)

		Note that on most terminals italic is rendered as underlining and
		depending on how old the terminal, other formatting might not appear
		as expected. If you know how to set LESS_TERMCAP_* variables they
		will be observed when output is to the terminal.

		Blocks

		1. Paragraph
		2. Verbatim (block begins with '    ', never first)
		3. Numbered (block begins with '* ')
		4. Bulleted (block begins with '1. ')

		Currently, a verbatim block must never be first because of the
		stripping of initial white space.

		Templates

		Anything from Cmd that fulfills the requirement to be included in
		a Go text/template may be used. This includes {{ "{{ cmd .Name }}" }}
		and the rest. A number of builtin template functions have also been
		added (such as {{ "indent" }}) which can receive piped input. You
		can add your own functions (or overwrite existing ones) by adding
		your own Dynamic template.FuncMap (see text/template for more about
		Go templates). Note that verbatim blocks will need to indented to work:

		    {{ "{{ dir | indent 4 }}" }}


		Produces a nice verbatim block:

		    {{ dir | indent 4 }}


		Note this is different for every user and their specific system. The
		ability to incorporate dynamic data into any help documentation is
		a game-changer not only for creating very consumable tools, but
		creating intelligent, interactive training and education materials
		as well.

		Templates Within Templates

		Sometimes you will need more text than can easily fit within
		a single action. (Actions may not span new lines.) For such things
		defining a template with that text is required and they you can
		include it with the {{pre "template" }} tag.

		    {{define "long" -}}
		    Here is something
		    that spans multiple
		    lines that would otherwise be too long for a single action.
		    {{- end}}

		Something
		`,
	/*
		Other: []Z.Section{
			{`Custom Sections`, `
				Additional sections can be added to the Other field.

				A Z.Section is just a Title and Body and can be assigned using
				composite notation (without the key names) for cleaner, in-code
				documentation.

				The Title will be capitalized for terminal output if using the
				common help.Cmd, but should use a suitable case for appearing in
				a book for other output renderers later (HTML, PDF, etc.)`,
			},
		},
	*/

	// no Call since has Commands, if had Call would only call if
	// commands didn't match
}

var GitAuthorCommitsCmd = &Z.Cmd{
	Name:    `authorcommits`,
	Summary: `lists the number of commits per author in current dir`,
	Aliases: []string{"ac"},
	Call: func(_ *Z.Cmd, _ ...string) error { // note conventional _

		w := new(tabwriter.Writer)

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()

		fmt.Fprintf(w, " %s\t%s\n", "Author", "Commits")
		fmt.Fprintf(w, " %s\t%s\n", "------", "-------")
		for k, v := range GitAuthorCommits() {
			fmt.Fprintf(w, " %s\t%d\n", k, v)
		}

		return nil
	},
	Commands: []*Z.Cmd{help.Cmd},
}

var GitMapAuthorChangesCmd = &Z.Cmd{
	Name:    `authorchanges`,
	Summary: `lists the line changes per author in current branch`,
	Aliases: []string{"ach"},
	Call: func(_ *Z.Cmd, _ ...string) error { // note conventional _

		w := new(tabwriter.Writer)

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)
		defer w.Flush()

		fmt.Fprintf(w, " %s\t%s\t%s\n", "Author", "Additions", "Deletions")
		fmt.Fprintf(w, " %s\t%s\t%s\n", "------", "---------", "---------")
		for k, v := range MapLineChanges() {
			fmt.Fprintf(w, " %s\t%d\t%d\n", k, v.Additions, v.Deletions)
		}

		return nil
	},
	Commands: []*Z.Cmd{help.Cmd},
}

var ContributionSummaryCmd = &Z.Cmd{
	Name:    `summary`,
	Summary: `lists the line changes per author in current branch`,
	Aliases: []string{"s"},
	Call: func(_ *Z.Cmd, _ ...string) error { // note conventional _

		commitMap := GitAuthorCommits()
		lineChangesMap := MapLineChanges()
		commitRatioMap := make(map[string]float64)
		lineRatioMap := make(map[string]float64)
		granularityMap := make(map[string]float64)

		// calculate aggregate metrics
		commitTotal := 0
		for _, v := range commitMap {
			commitTotal += v
		}

		lineTotal := 0
		for _, v := range lineChangesMap {
			lineTotal += v.Sum()
		}

		for k, v := range lineChangesMap {
			linesum := v.Sum()
			lineRatio := float64(linesum) / float64(lineTotal)
			lineRatioMap[k] = lineRatio
			commitRatio := float64(commitMap[k]) / float64(commitTotal)
			commitRatioMap[k] = commitRatio
			granularityMap[k] = 1.0 / (float64(linesum) / float64(commitMap[k]))
		}

		// output results
		w := new(tabwriter.Writer)

		// minwidth, tabwidth, padding, padchar, flags
		w.Init(os.Stdout, 8, 8, 0, '\t', 0)

		fmt.Fprintf(w, " %s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Author", "Commits", "Additions", "Deletions", "Line ratio", "Commit ratio", "Granularity")
		fmt.Fprintf(w, " %s\t%s\t%s\t%s\t%s\t%s\t%s\n", "------", "-------", "---------", "---------", "----------", "------------", "-----------")
		for k, v := range MapLineChanges() {
			fmt.Fprintf(w, " %s\t%v\t%v\t%v\t%.2f\t%.2f\t%.2f\n", k, commitMap[k], v.Additions, v.Deletions, lineRatioMap[k], commitRatioMap[k], granularityMap[k])
		}
		err := w.Flush()
		if err != nil {
			return fmt.Errorf("failed to flush output buffer: %w", err)
		}

		fmt.Printf(
			"\n Overall repo commit granularity: %.2f\n",
			1.0/(float64(lineTotal)/float64(commitTotal)),
		)

		return nil
	},
	Commands: []*Z.Cmd{help.Cmd},
}
