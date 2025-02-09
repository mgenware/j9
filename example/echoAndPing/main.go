package main

import "github.com/mgenware/j9/v3"

func main() {
	t := j9.NewTunnel(j9.NewLocalNode(), j9.NewConsoleLogger())
	t.Spawn(&j9.SpawnParams{Name: "echo", Args: []string{"param1", "param2", "param3"}})
	t.Spawn(&j9.SpawnParams{Name: "ping", Args: []string{"google.com"}})
}
