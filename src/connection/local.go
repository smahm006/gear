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
	for k, v := range env_map {
		env = append(env, "%s=%s", k, v)
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
	cmd = exec.Command("uname -r")
	out, err := cmd.Output()
	if err != nil {
		return err
	} else {
		l.Host.SetOs(string(out))
	}
	cmd = exec.Command("cat /etc/*release")
	out, err = cmd.Output()
	if err != nil {
		return err
	} else {
		l.Host.SetDistro(string(out))
	}
	cmd.Env = EnvMapToSlice(l.Host.Environment)
	cmd = exec.Command("")
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

func (l *LocalConnection) Execute() error {
	stdout_buffer := new(strings.Builder)
	stderr_buffer := new(strings.Builder)
	l.Session.Stdout = stdout_buffer
	l.Session.Stderr = stderr_buffer
	sh := "env sh"
	l.Session.Path = sh
	l.Session.Args = []string{sh, "-c"}
	err := l.Session.Run()
	l.Session.Stdout = nil
	l.Session.Stderr = nil
	l.Session.Process = nil
	l.Session.ProcessState = nil
	return nil
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
