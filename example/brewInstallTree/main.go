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
		t.Spawn(&j9.SpawnParams{Name: "brew", Args: []string{"install", "tree"}})
	}
	fmt.Println("tree is installed")
	t.Spawn(&j9.SpawnParams{Name: "tree", Args: []string{"."}})
}
