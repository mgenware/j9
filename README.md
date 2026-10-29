# j9

[![Build](https://github.com/mgenware/j9/actions/workflows/build.yml/badge.svg)](https://github.com/mgenware/j9/actions/workflows/build.yml)
[![Lint](https://github.com/mgenware/j9/actions/workflows/lint.yml/badge.svg)](https://github.com/mgenware/j9/actions/workflows/lint.yml)

Shell scripting, SSH in Go.

## Installation

```sh
go get github.com/mgenware/j9/v3
```

## Examples

Check if `tree` command is available and install it if necessary. (Assuming on macOS with homebrew installed):

```go
package main

import (
	"fmt"
	"os/exec"

	"github.com/mgenware/j9/v3"
)

func main() {
	// Local node runs on your local system.
	// `ConsoleLogger` prints logs to the current console.
	t := j9.NewTunnel(j9.NewLocalNode(), j9.NewConsoleLogger())

	// Check if the command `tree` is installed.
	_, err := exec.LookPath("tree")
	if err != nil {
		t.Logger().Log(j9.LogLevelError, "tree is not installed")
		t.Spawn(&j9.SpawnOpt{Name: "brew", Args: []string{"install", "tree"}})
	}
	fmt.Println("tree is installed")
	t.Spawn(&j9.SpawnOpt{Name: "tree", Args: []string{"."}})
}
```

Sample output when `tree` is not installed:

```
❯ go run main.go
tree is not installed
brew install tree
==> Downloading https://homebrew.bintray.com/bottles/tree-1.7.0.high_sierra.bottle.1.tar.gz
==> Pouring tree-1.7.0.high_sierra.bottle.1.tar.gz
🍺  /usr/local/Cellar/tree/1.7.0: 8 files, 114.3KB
tree is installed
tree .
.
└── main.go

0 directories, 1 file
```

### SSH with a private key

SSH server host keys are verified against `~/.ssh/known_hosts` by default. Ensure the server's key has been verified and added to that file before connecting. To use a different verification policy, set `SSHConfig.HostKeyCallback`.

```go
package main

import (
	"github.com/mgenware/j9/v3"
)

func main() {
	config := &j9.SSHConfig{
		Host: "1.2.3.4",
		User: "root",
		Auth: j9.MustCreateKeyBasedAuth("~/key.pem"),
	}

	t := j9.NewTunnel(j9.MustCreateSSHNode(config), j9.NewConsoleLogger())
	t.Shell("pwd")
	t.Shell("ls")
}
```

Sample output:

```
pwd
/root

ls
bin
build
data
```

## Windows Support

Use WSL 2.

## Spawn vs Shell

Use them based on your use cases.

|                                                       | Spawn | Shell |
| ----------------------------------------------------- | ----- | ----- |
| Return stdout and stderr as a string                  | ❌    | ✅    |
| Live process output (good for long-running processes) | ✅    | ❌    |
| Supported in `LocalNode`                              | ✅    | ✅    |
| Supported in `SSHNode`                                | ❌    | ✅    |
