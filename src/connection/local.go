package connection

import (
	"fmt"
	"os"
	"os/exec"
)

type LocalConnection struct {
	Environment []string
	Session     *exec.Cmd
}

func NewLocalConnection(localhost_env []string) *LocalConnection {
	return &LocalConnection{Environment: localhost_env}
}

func (l *LocalConnection) Connect() error {
	if l.Session != nil {
		return nil
	}
	cmd := exec.Command("/bin/bash")
	cmd.Env = l.Environment
	l.Session = cmd
	return nil
}

func (l *LocalConnection) Close() error {
	return nil
}

func (l *LocalConnection) WhoAmI() (string, error) {
	user, exists := os.LookupEnv("USER")
	if !exists {
		return user, fmt.Errorf("could not read environment variable $USER")
	}
	return user, nil
}

func (l *LocalConnection) Execute() error {
	return nil
}

func (l *LocalConnection) CopyFile() error {
	return nil
}

func (l *LocalConnection) WriteData() error {
	return nil
}

// func (l *LocalConnection) Run(command string) {
// 	l.Session.Path = "/bin/bash"
// 	l.Session.Args = []string{"/bin/bash", "-c", command}
// 	out, _ := l.Session.CombinedOutput()
// 	fmt.Print("OUT: ", string(out))
// 	l.Session.Stdout = nil
// 	l.Session.Stderr = nil
// 	l.Session.Process = nil
// 	l.Session.ProcessState = nil
// }
