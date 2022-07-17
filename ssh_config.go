package j9

import (
	"io/ioutil"

	"github.com/mgenware/j9/lib"
	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host string
	User string
	Port int
	Auth []ssh.AuthMethod
}

func SafeNewKeyBasedAuth(keyFile string) ([]ssh.AuthMethod, error) {
	keyBytes, err := ioutil.ReadFile(lib.FormatPath(keyFile, true))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		return nil, err
	}
	return []ssh.AuthMethod{ssh.PublicKeys(signer)}, nil
}

func NewKeyBasedAuth(keyFile string) []ssh.AuthMethod {
	auth, err := SafeNewKeyBasedAuth(keyFile)
	if err != nil {
		panic(err)
	}
	return auth
}

func NewDefaultKeyBasedAuth() []ssh.AuthMethod {
	return NewKeyBasedAuth("~/.ssh/id_rsa")
}

func SafeNewPwdBasedAuth(pwd string) ([]ssh.AuthMethod, error) {
	return []ssh.AuthMethod{ssh.Password(pwd)}, nil
}

func NewPwdBasedAuth(pwd string) []ssh.AuthMethod {
	auth, err := SafeNewPwdBasedAuth(pwd)
	if err != nil {
		panic(err)
	}
	return auth
}
