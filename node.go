package j9

// Node is an interface for running commands in a specific environment.
type Node interface {
	// RunCmd runs the given command, returns an error if the command fails.
	RunCmd(wd string, name string, arg ...string) error

	// RunCmdSync runs the given command, returns the output and an error if the command fails.
	RunCmdSync(wd string, cmd string) ([]byte, error)
}
