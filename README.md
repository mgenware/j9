# j9

[![Build](https://github.com/mgenware/j9/actions/workflows/build.yml/badge.svg)](https://github.com/mgenware/j9/actions/workflows/build.yml)
[![Lint](https://github.com/mgenware/j9/actions/workflows/lint.yml/badge.svg)](https://github.com/mgenware/j9/actions/workflows/lint.yml)

Shell scripting in Go.

## Installation

```sh
go get github.com/mgenware/j9
```

## Examples

Check if `tree` command is available and install it if necessary. (Assuming on macOS with homebrew installed):

```go
package main

import (
	"fmt"
	"os/exec"

	"github.com/mgenware/j9/v2"
)

func main() {
	// Local node runs on your local system.
	// `ConsoleLogger` prints logs to the current console.
	t := j9.NewTunnel(j9.NewLocalNode(), j9.NewConsoleLogger())

	// Check if the command `tree` is installed.
	_, err := exec.LookPath("tree")
	if err != nil {
		t.Logger().Log(j9.LogLevelError, "tree is not installed")
		t.Run("brew", "install", "tree")
	}
	fmt.Println("tree is installed")
	t.Run("tree", ".")
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

```go
package main

import (
	"github.com/mgenware/j9/v2"
)

func main() {
	config := &j9.SSHConfig{
		Host: "1.2.3.4",
		User: "root",
		Auth: j9.MustCreateKeyBasedAuth("~/key.pem"),
	}

	t := j9.NewTunnel(j9.MustCreateSSHNode(config), loggers.NewConsoleLogger())
	t.Run("pwd")
	t.Run("ls")
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
