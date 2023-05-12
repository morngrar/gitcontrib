// Copyright 2023 gitcontrib Authors
// SPDX-License-Identifier: Apache-2.0

// Package example provides the Bonzai command branch of the same name.
package gitcontrib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"text/template"

	Z "github.com/rwxrob/bonzai/z"
	"github.com/rwxrob/help"
)

// Cmd provides a Bonzai branch command that can be composed into Bonzai
// trees or used as a standalone with light wrapper (see cmd/).
var Cmd = &Z.Cmd{

	Name:      `gitcontrib`,
	Summary:   `a command tree branch for analysing git author contributions`,
	Version:   `v0.1.1`,
	Copyright: `Copyright 2023 Svein-Kåre Bjørnsen`,
	License:   `Apache-2.0`,
	Source:    `git@git.gvk.idi.ntnu.no:morngrar/gitcontrib.git`,
	Issues:    `https://git.gvk.idi.ntnu.no/morngrar/gitcontrib/-/issues`,

	// Composite commands, local and external, all have their own names
	// that are added to the command tree depending on where they are
	// composed.

	Commands: []*Z.Cmd{

		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd,

		// local commands (in this module)
		AuthorCommitsCmd, AuthorChangesCmd, ContributionSummaryCmd,
		CsvCmd,
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
		The {{aka}} command is a command tree for analysing author
		contributions in git repositories. The most extensive overview is the
		'summary' subcommand, that has the alias 's'. This shows aggregated
		metrics like line change ratio and commit granularity per author.
		Larger numbers means higher contribution for all categories.
		`,
}

var AuthorCommitsCmd = &Z.Cmd{
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
		for k, v := range AuthorCommits() {
			fmt.Fprintf(w, " %s\t%d\n", k, v)
		}

		return nil
	},
	Commands: []*Z.Cmd{help.Cmd},
}

var AuthorChangesCmd = &Z.Cmd{
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
	Summary: `lists commits, line changes and aggregated metrics`,
	Aliases: []string{"s"},
	Call: func(_ *Z.Cmd, _ ...string) error { // note conventional _

		commitMap := AuthorCommits()
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

		w.Init(os.Stdout, 8, 8, 0, '\t', 0) // setting up table dimensions

		fmt.Fprintf(w, " %s\t%s\t%s\t%s\t%s\t%s\t%s\n", "Author", "Commits", "Additions", "Deletions", "Line ratio", "Commit ratio", "Granularity")
		fmt.Fprintf(w, " %s\t%s\t%s\t%s\t%s\t%s\t%s\n", "------", "-------", "---------", "---------", "----------", "------------", "-----------")
		for k, v := range MapLineChanges() {
			fmt.Fprintf(w, " %s\t%v\t%v\t%v\t%.3f\t%.3f\t%.3f\n", k, commitMap[k], v.Additions, v.Deletions, lineRatioMap[k], commitRatioMap[k], granularityMap[k])
		}
		err := w.Flush()
		if err != nil {
			return fmt.Errorf("failed to flush output buffer: %w", err)
		}

		fmt.Printf(
			"\n Overall repo commit granularity: %.3f\n",
			1.0/(float64(lineTotal)/float64(commitTotal)),
		)

		return nil
	},
	Commands: []*Z.Cmd{help.Cmd},
}

var CsvCmd = &Z.Cmd{
	Name:    `csv`,
	Summary: `outputs CSV rows for the various report`,
	Aliases: []string{"c"},
	Commands: []*Z.Cmd{

		// standard external branch imports (see rwxrob/{help,conf,vars})
		help.Cmd,

		// local commands (in this module)
		CsvContributionSummaryCmd,
	},
	Description: `
		The {{aka}} subcommand supplies the same commands as the root command, 
		however, these output CSV rows instead of the human-readable tabulated
		output of the original commands. The first field of each row is the
		name of the repo directory itself, the rest follow the same order as 
		the original command. Strings are wrapped in double quotes.
		`,
}

var CsvContributionSummaryCmd = &Z.Cmd{
	Name:    `summary`,
	Summary: `outputs CSV rows for the 'summary' report`,
	Aliases: []string{"s"},
	Description: `
		The {{aka}} subcommand gives the same output data as the  root 
		subcommand of the same name, 
		however, this one outputs CSV rows instead of the human-readable tabulated
		output of the original command. The first field of each row is the
		name of the repo directory itself, the rest follow the same order as 
		the original command. Strings are wrapped in double quotes, and the CSV 
		header is not printed to accomodate scripting.

		The fields of this command is the following, in the given order:

		Repo directory, Author, Commits, Additions, Deletions, Line ratio, Commit ratio, Granularity
		`,

	Call: func(_ *Z.Cmd, _ ...string) error { // note conventional _

		commitMap := AuthorCommits()
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

		reponame, err := getRepoDirName()
		if err != nil {
			return fmt.Errorf("error getting repo name: %w", err)
		}

		for k, v := range MapLineChanges() {
			fmt.Printf("\"%s\",\"%s\",%v,%v,%v,%.3f,%.3f,%.3f\n", reponame, k, commitMap[k], v.Additions, v.Deletions, lineRatioMap[k], commitRatioMap[k], granularityMap[k])
		}

		return nil
	},
	Commands: []*Z.Cmd{help.Cmd},
}

func getRepoDirName() (string, error) {

	output := Z.Out("git", "rev-parse", "--show-toplevel")
	if output == "" {
		return "", errors.New("error getting git repo directory path")
	}

	dirname := strings.TrimSpace(filepath.Base(output))

	return dirname, nil
}
