package j9

type Node interface {
	RunOrError(cmd string) ([]byte, error)
}
