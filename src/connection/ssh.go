package connection

import (
	"github.com/smahm006/gear/src/inventory"
	"golang.org/x/crypto/ssh"
)

type SshConnection struct {
	Hostname      string
	Username      string
	Port          int
	Environment   []string
	Session       *ssh.Session
	ForwardAgent  bool
	LoginPassword *string
	Key           *string
	Passphrase    *string
	KeyComment    *string
}

func NewSshConnection(host *inventory.Host, remote_env []string) *SshConnection {
	return &SshConnection{
		Hostname:    host.Name,
		Environment: remote_env,
	}
}

func (l *SshConnection) Connect() error {
	return nil
}

func (l *SshConnection) Close() error {
	return nil
}

func (l *SshConnection) WhoAmI() (string, error) {
	return "", nil
}

func (l *SshConnection) Execute() error {
	return nil
}

func (l *SshConnection) CopyFile() error {
	return nil
}

func (l *SshConnection) WriteData() error {
	return nil
}
