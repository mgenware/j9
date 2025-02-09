package j9

type SpawnParams struct {
	WorkingDir string
	Name       string
	Args       []string
	Env        []string
}

type ShellParams struct {
	WorkingDir string
	Cmd        string
	Env        []string
}

// Node is an interface for running commands in a specific environment.
type Node interface {
	Spawn(params *SpawnParams) error
	Shell(params *ShellParams) (string, error)
}
