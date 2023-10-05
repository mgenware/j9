package main

import "github.com/mgenware/j9/v2"

func main() {
	ln := j9.NewLocalNode()
	lt := j9.NewTunnel(ln, j9.NewConsoleLogger())
	lt.Run("echo", "param1", "param2", "param3")
	lt.Run("ping", "google.com")
}
