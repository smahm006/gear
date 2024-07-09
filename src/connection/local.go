package connection

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/smahm006/gear/src/inventory"
	"github.com/smahm006/gear/src/utils"
)

type LocalConnection struct {
	Host    *inventory.Host
	Session *exec.Cmd
}

func EnvMapToSlice(env_map map[string]string) []string {
	var env []string
	env = append(env, "PATH=/usr/sbin:/usr/bin:/sbin:/bin")
	for k, v := range env_map {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	return env
}

func NewLocalConnection(Host *inventory.Host) *LocalConnection {
	return &LocalConnection{Host: Host}
}

func (l *LocalConnection) Connect() error {
	var cmd *exec.Cmd
	if l.Session != nil {
		return nil
	}
	cmd = exec.Command("env", "sh", "-c", "uname -r")
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return err
	} else {
		l.Host.SetOs(string(out))
	}
	cmd = exec.Command("env", "sh", "-c", "cat /etc/*release")
	out, err = cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return err
	} else {
		l.Host.SetDistro(string(out))
	}
	cmd = exec.Command("")
	cmd.Env = EnvMapToSlice(l.Host.Environment)
	l.Session = cmd
	return nil
}

func (l *LocalConnection) Close() error {
	err := l.Session.Run()
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalConnection) WhoAmI() (string, error) {
	user, exists := os.LookupEnv("USER")
	if !exists {
		return user, fmt.Errorf("could not read environment variable $USER")
	}
	return user, nil
}

func (l *LocalConnection) Execute(command string) (string, error) {
	stdout_buffer := new(strings.Builder)
	stderr_buffer := new(strings.Builder)
	l.Session.Stdout = stdout_buffer
	l.Session.Stderr = stderr_buffer
	l.Session.Path = "/bin/env"
	l.Session.Args = []string{"/bin/env", "sh", "-c", command}
	err := l.Session.Run()
	if err != nil {
		return "", err
	}
	l.Session.Stdout = nil
	l.Session.Stderr = nil
	l.Session.Process = nil
	l.Session.ProcessState = nil
	return stdout_buffer.String(), nil
}

func (l *LocalConnection) CopyFile(src string, dst string) error {
	err := utils.CopyFile(src, dst)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalConnection) WriteData(data string, path string) error {
	var file *os.File
	var err error
	file, err = utils.OpenFile(path)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}
	return nil
}
