package j9

import (
	"os"
	"strings"

	"github.com/mgenware/goutil/stringsx"
	"golang.org/x/crypto/ssh"
)

const homeEnv = "$HOME"

// SSHConfig is the configuration for SSHNode.
type SSHConfig struct {
	Host string
	User string
	Port int
	Auth []ssh.AuthMethod
}

// Creates a new ssh.AuthMethod with the given key file.
func NewKeyBasedAuth(keyFile string) ([]ssh.AuthMethod, error) {
	keyBytes, err := os.ReadFile(formatPath(keyFile, true))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
}

// Creates a new ssh.AuthMethod with the given key file, panics if there is an error.
func MustCreateKeyBasedAuth(keyFile string) []ssh.AuthMethod {
	auth, err := NewKeyBasedAuth(keyFile)
	if err != nil {
		panic(err)
	}
	return auth
}

// Creates a new ssh.AuthMethod with the default key file("~/.ssh/id_rsa").
func MustCreateDefaultKeyBasedAuth() []ssh.AuthMethod {
	return MustCreateKeyBasedAuth("~/.ssh/id_rsa")
}

// Creates a new ssh.AuthMethod with the given password.
func NewPwdBasedAuth(pwd string) ([]ssh.AuthMethod, error) {
	return []ssh.AuthMethod{ssh.Password(pwd)}, nil
}

// Creates a new ssh.AuthMethod with the given password, panics if there is an error.
func MustCreateNewPwdBasedAuth(pwd string) []ssh.AuthMethod {
	auth, err := NewPwdBasedAuth(pwd)
	if err != nil {
		panic(err)
	}
	return auth
}

func formatPath(s string, evaluate bool) string {
	if strings.HasPrefix(s, "~/") {
		s = homeEnv + "/" + stringsx.SubStringFromStart(s, 2)
	} else if s == "~" {
		s = homeEnv
	}
	if evaluate {
		return os.ExpandEnv(s)
	}
	return s
}
