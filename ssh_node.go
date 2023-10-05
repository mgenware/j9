package j9

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// SSHNode is a node that runs commands on a remote SSH server.
type SSHNode struct {
	dir       *dirManager
	config    *SSHConfig
	sshConfig *ssh.ClientConfig
	client    *ssh.Client

	Logger Logger
}

// Creates a new SSH node with the given configuration.
func NewSSHNode(config *SSHConfig) (*SSHNode, error) {
	if config.Host == "" {
		return nil, fmt.Errorf("config.Host cannot be empty")
	}
	if config.User == "" {
		return nil, fmt.Errorf("config.User cannot be empty")
	}
	if len(config.Auth) == 0 {
		return nil, fmt.Errorf("config.Auth cannot be empty")
	}
	if config.Port == 0 {
		config.Port = 22
	}
	sshConfig := &ssh.ClientConfig{
		User:            config.User,
		Auth:            config.Auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &SSHNode{
		config:    config,
		sshConfig: sshConfig,
		dir:       &dirManager{},
	}, nil
}

// Creates a new SSH node with the given configuration, panics if there is an error.
func MustCreateSSHNode(config *SSHConfig) *SSHNode {
	node, err := NewSSHNode(config)
	if err != nil {
		panic(err)
	}
	return node
}

func (node *SSHNode) RunOrError(cmd string) ([]byte, error) {
	session, err := node.prepareSession()
	if err != nil {
		return nil, err
	}
	output, err := node.runCore(session, cmd)
	if err != nil {
		if len(output) > 0 {
			node.log(LogLevelWarning, string(output))
		}
		node.log(LogLevelWarning, err.Error())
		node.log(LogLevelWarning, "Reconnecting...")

		// set the previous client to nil
		if node.client != nil {
			err = node.client.Close()
			if err != nil {
				node.log(LogLevelWarning, "Error closing previous client: "+err.Error())
			}
			node.client = nil
		}
		session, err = node.prepareSession()
		if err != nil {
			return nil, err
		}

		output, err = node.runCore(session, cmd)
		if err != nil {
			return output, err
		}
	}
	return output, nil
}

func (node *SSHNode) runCore(session *ssh.Session, cmd string) ([]byte, error) {
	node.dir.NextWD(cmd, false)

	lastDir := node.dir.LastDir()
	if lastDir != "" {
		// TODO: better handling of escaping path
		cmd = "cd '" + lastDir + "' && " + cmd
	}

	return session.CombinedOutput(cmd)
}

func (node *SSHNode) prepareSession() (*ssh.Session, error) {
	if node.client == nil {
		cfg := node.config
		addr := fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
		node.log(LogLevelInfo, fmt.Sprintf("SSH: Connecting to %v\n", addr))

		client, err := ssh.Dial("tcp", addr, node.sshConfig)
		if err != nil {
			return nil, err
		}
		node.client = client
	}

	session, err := node.client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (node *SSHNode) log(logLevel int, message string) {
	if node.Logger != nil {
		node.Logger.Log(logLevel, message)
	}
}
