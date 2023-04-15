// Copyright 2022 gitcontrib Authors
// SPDX-License-Identifier: Apache-2.0

package gitcontrib

import (
	"fmt"

	Z "github.com/rwxrob/bonzai/z"
)

func GitAuthorCommits() string {
	var out string

	out = Z.Out("git", "branch")
	fmt.Println(out)
	// TODO: determine main branch (master/main) and use it for the below command

	// git branch has to be passed when invoking like this
	// https://stackoverflow.com/questions/51966053/what-is-wrong-with-invoking-git-shortlog-from-go-exec
	out = Z.Out("git", "shortlog", "-sn", "--no-merges", "main")
	return out
}

func MapAddedAndRemoved() {

	// TODO: perform something like this:

	//	#!/bin/sh
	//	declare -A map
	//	while read line; do
	//		if grep "^[a-zA-Z]" <<< "$line" > /dev/null; then
	//			current="$line"
	//			if [ -z "${map[$current]}" ]; then
	//				map[$current]=0
	//			fi
	//		elif grep "^[0-9]" <<<"$line" >/dev/null; then
	//			for i in $(cut -f 1,2 <<< "$line"); do
	//				map[$current]=$((map[$current] + $i))
	//			done
	//		fi
	//	done <<< "$(git log --numstat --pretty="%aN")"
	//
	//	for i in "${!map[@]}"; do
	//		echo -e "$i:${map[$i]}"
	//	done | sort -nr -t ":" -k 2 | column -t -s ":"

}

func ContributionStatus() {
	// TODO: combine the two functions above into a comprehensive report
}
