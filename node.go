package j9

// Node is an interface for running commands in a specific environment.
type Node interface {
	// RunUnsafe runs the given command, returns an error if the command fails.
	RunUnsafe(name string, arg ...string) error

	// RunSyncUnsafe runs the given command, returns the output and an error if the command fails.
	RunSyncUnsafe(cmd string) ([]byte, error)

	CDUnsafe(dir string) error
}
