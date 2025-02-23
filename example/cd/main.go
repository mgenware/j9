package main

import (
	"fmt"

	"github.com/mgenware/j9/v3"
)

func main() {
	t := j9.NewTunnel(j9.NewLocalNode(), j9.NewConsoleLogger())
	fmt.Println("> Current dir")

	t.Spawn(&j9.SpawnOpt{Name: "pwd"})

	fmt.Println("> Parent dir")
	t.CD("..")
	t.Spawn(&j9.SpawnOpt{Name: "pwd"})

	fmt.Println("> CD example")
	t.CD("example")
	t.Spawn(&j9.SpawnOpt{Name: "pwd"})
}
