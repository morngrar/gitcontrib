// Copyright 2022 gitcontrib Authors
// SPDX-License-Identifier: Apache-2.0

package gitcontrib

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	Z "github.com/rwxrob/bonzai/z"
)

func extractCheckedOutBranch(gitBranchOutput string) (string, error) {

	// for all lines in git branch output, find the active one
	scanner := bufio.NewScanner(strings.NewReader(gitBranchOutput))
	var branch string
	for scanner.Scan() {
		line := scanner.Text()
		match, err := regexp.MatchString("^\\* *", line)
		if err != nil {
			return "", fmt.Errorf("error matching line: %w", err)
		}

		if match {
			branch = strings.Fields(line)[1]
			break
		}
	}

	return branch, nil
}

func mapAuthorCommits(shortlogOutput string) (map[string]int, error) {

	authorMap := make(map[string]int)
	scanner := bufio.NewScanner(strings.NewReader(shortlogOutput))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		commits, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("error parsing commit number: %w", err)
		}

		authorMap[strings.Join(fields[1:], " ")] = commits

	}

	return authorMap, nil
}

// GitAuthorCommits returns a map of author names with their respective
// non-merge commit counts as values
func GitAuthorCommits() map[string]int {
	var out string

	out = Z.Out("git", "branch")
	branch, err := extractCheckedOutBranch(out)
	if err != nil {
		log.Fatalf("Error extracting branch: %s", err)
	}

	// git branch has to be passed when invoking like this
	// https://stackoverflow.com/questions/51966053/what-is-wrong-with-invoking-git-shortlog-from-go-exec
	out = Z.Out("git", "shortlog", "-sn", "--no-merges", branch)
	authorMap, err := mapAuthorCommits(out)
	if err != nil {
		log.Fatalf("Error extracting commit counts: %s", err)
	}

	return authorMap
}

type LineChanges struct {
	Additions int
	Deletions int
}

func (lc *LineChanges) Add(n int) {
	lc.Additions += n
}

func (lc *LineChanges) Del(n int) {
	lc.Deletions += n
}

func (lc *LineChanges) Sum() int {
	return lc.Additions + lc.Deletions
}

func parseLineChanges(gitOutput string) (map[string]LineChanges, error) {
	authorMap := make(map[string]LineChanges)

	scanner := bufio.NewScanner(strings.NewReader(gitOutput))
	currentAuthor := ""
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// check if line is author
		match, err := regexp.MatchString("^[a-zA-Z']", line)
		if match {
			if strings.HasPrefix(line, "'") && strings.HasSuffix(line, "'") {
				line = line[:len(line)-1]
				line = line[1:]
			}
			currentAuthor = line
			_, ok := authorMap[line]
			if !ok { // new author
				authorMap[currentAuthor] = LineChanges{0, 0}
			}
			continue
		}

		// if not new author, accumulate counts
		var adds int
		var dels int

		fields := strings.Fields(line)
		if fields[0] != "-" {
			adds, err = strconv.Atoi(fields[0])
			if err != nil {
				return nil, fmt.Errorf("error parsing adds: %s", err)
			}
		}

		if fields[1] != "-" {
			dels, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, fmt.Errorf("error parsing dels: %s", err)
			}
		}

		a := authorMap[currentAuthor]
		a.Add(adds)
		a.Del(dels)
		authorMap[currentAuthor] = a
	}

	return authorMap, nil
}

// MapLineChanges returns an author map containing the line changes of each
// author in the current repo branch.
func MapLineChanges() map[string]LineChanges {

	out := Z.Out("git", "log", "--numstat", "--pretty='%aN'")
	authorMap, err := parseLineChanges(out)
	if err != nil {
		log.Fatalf("Error extracting commit counts: %s", err)
	}

	return authorMap
}
