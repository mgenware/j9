package main

import (
	"fmt"

	"github.com/mgenware/j9/v2"
)

func main() {
	ln := j9.NewLocalNode()
	lt := j9.NewTunnel(ln, j9.NewConsoleLogger())
	fmt.Println("> Current dir")
	lt.Run("pwd")

	fmt.Println("> Parent dir")
	lt.CD("..")
	lt.Run("pwd")

	fmt.Println("> CD example")
	lt.CD("example")
	lt.Run("pwd")
}
