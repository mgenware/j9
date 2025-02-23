package j9

type SpawnOpt struct {
	WorkingDir string
	Name       string
	Args       []string
	Env        []string
}

type ShellOpt struct {
	WorkingDir string
	Cmd        string
	Env        []string
}

// Node is an interface for running commands in a specific environment.
type Node interface {
	Spawn(params *SpawnOpt) error
	Shell(params *ShellOpt) (string, error)
}
