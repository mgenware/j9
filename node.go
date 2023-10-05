package j9

// Node is an interface for running commands on a specific environment.
type Node interface {
	RunOrError(cmd string) ([]byte, error)
}
