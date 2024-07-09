package connection

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/tasks/requonse"
	"github.com/smahm006/gear/internal/utils"
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
	cmd = exec.Command("env", "sh", "-c", "uname")
	out, err := cmd.Output()
	if err != nil {
		return err
	} else {
		l.Host.SetOs(string(out))
	}
	cmd = exec.Command("env", "sh", "-c", "cat /etc/*release")
	out, err = cmd.Output()
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

func (l *LocalConnection) Execute(command string) *requonse.TaskResponse {
	var outbuf, errbuf strings.Builder
	var exitcode int
	response := requonse.NewTaskResponse()
	l.Session.Path = "/usr/bin/env"
	l.Session.Stdout = &outbuf
	l.Session.Stderr = &errbuf
	l.Session.Args = []string{"/usr/bin/env", "sh", "-c", fmt.Sprintf("LANG=C %s", command)}
	err := l.Session.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()
	if err != nil {
		response.Type = requonse.Failed
		if exit_error, ok := err.(*exec.ExitError); ok {
			ws := exit_error.Sys().(syscall.WaitStatus)
			exitcode = ws.ExitStatus()
		} else {
			exitcode = default_exit_code
			if stderr == "" {
				stderr = fmt.Sprintf(err.Error())
			}
		}
	} else {
		ws := l.Session.ProcessState.Sys().(syscall.WaitStatus)
		exitcode = ws.ExitStatus()
	}
	response.CommandResult = &requonse.CommandResult{
		Cmd: command,
		Out: stdout,
		Err: stderr,
		Rc:  exitcode,
	}
	l.Session.Stdout = nil
	l.Session.Stderr = nil
	l.Session.Process = nil
	l.Session.ProcessState = nil
	return response
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
