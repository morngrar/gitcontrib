# gitcontrib: a Bonzai™ branch and standalone tool for git analysis


[![License](https://img.shields.io/badge/license-Apache2-brightgreen.svg)](LICENSE)

## Usage

**TODO**

## Install

This command can be installed as a standalone program or composed into a
Bonzai command tree.

Standalone

```
go install github.com/PLACEHOLDER/gitcontrib/cmd/gitcontrib@latest
```

Composed into a Bonzai™ command tree

```go
package z

import (
	Z "github.com/rwxrob/bonzai/z"
	gitcontrib "github.com/PLACEHOLDER/gitcontrib"
)

var Cmd = &Z.Cmd{
	Name:     `z`,
	Commands: []*Z.Cmd{help.Cmd, gitcontrib.Cmd, gitcontrib.BazCmd},
}
```

## What is a Bonzai™ command tree?

Bonzai is a framework for creating composable command trees, which are
multiplatform single-executable toolkits. Each subcommand can be imported into
any other such command tree, or work as a standalone command tree, so that each
toolkit creator can mix and match branches in their own command trees. It also
comes with a simple way to autocomplete under bash and having embedded dynamic
documentation.

See [the library repo][https://github.com/rwxrob/bonzai] for more details, and
[the original command tree](https://github.com/rwxrob/z) that seeded the
creation of the framework for a concrete example other than gitcontrib.

## Tab Completion

To activate bash completion just use the `complete -C` option from your
`.bashrc` or command line. There is no messy sourcing required. All the
completion is done by the program itself.

```
complete -C gitcontrib gitcontrib
```

If you don't have bash or tab completion check use the shortcut
commands instead.

## Embedded Documentation

All documentation (like manual pages) has been embedded into the source
code of the application. See the source or run the program with help to
access it.

## Copyright notice

This repo used the [bonzai example](https://github.com/rwxrob/bonzai-example)
template as a starting point. It is heavily modified, but some fragments may
remain. Bonzai and that example is licensed under the Apache-2.0 license.

