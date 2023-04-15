// Copyright 2023 gitcontrib Authors
// SPDX-License-Identifier: Apache-2.0

package gitcontrib

import (
	Z "github.com/rwxrob/bonzai/z"
)

// private leaf
var gitStatusCmd = &Z.Cmd{
	Name: `gstatus`,
	Call: func(caller *Z.Cmd, none ...string) error {
		Z.Exec("git", "status")
		return nil
	},
}
