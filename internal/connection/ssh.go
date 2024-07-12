package connection

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pkg/sftp"
	"github.com/smahm006/gear/internal/inventory"
	"github.com/smahm006/gear/internal/tasks/exchange"
	"github.com/smahm006/gear/internal/utils"
	"golang.org/x/crypto/ssh"
)

type SshConnection struct {
	Host    *inventory.Host
	Config  *ssh.ClientConfig
	Client  *ssh.Client
	Session *ssh.Session
}

func NewSshConnection(Host *inventory.Host) *SshConnection {
	return &SshConnection{
		Config: &ssh.ClientConfig{},
		Host:   Host,
	}
}

func (s *SshConnection) Connect() error {
	// SSH username
	var username string
	if gear_ssh_user := s.Host.Getvar("gear_ssh_user"); gear_ssh_user != 0 {
		username = fmt.Sprint(gear_ssh_user)
	} else if env_gear_ssh_user := os.Getenv("GEAR_SSH_USER"); len(env_gear_ssh_user) != 0 {
		username = env_gear_ssh_user
	} else {
		env_u := os.Getenv("USER")
		username = env_u
	}
	s.Config.User = username
	// SSH Authentication
	// Always prefer key authentication over password if possible
	var signer ssh.Signer
	var password string
	var err error
	if gear_ssh_pkey := s.Host.Getvar("gear_ssh_private_key"); gear_ssh_pkey != nil {
		var pk_data []byte
		pk_data, err = utils.ReadFile(fmt.Sprint(gear_ssh_pkey))
		if err == nil {
			signer, err = ssh.ParsePrivateKey(pk_data)
		}
	} else if env_gear_ssh_pkey := os.Getenv("GEAR_SSH_PRIVATE_KEY"); err != nil && len(env_gear_ssh_pkey) != 0 {
		pk_data, err := utils.ReadFile(env_gear_ssh_pkey)
		if err == nil {
			signer, err = ssh.ParsePrivateKey(pk_data)
		}
	}
	if signer != nil {
		s.Config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		// Could not find or parse a private key
		if gear_ssh_password := s.Host.Getvar("gear_ssh_password"); gear_ssh_password != nil {
			password = fmt.Sprint(gear_ssh_password)
		} else if env_gear_ssh_password := os.Getenv("GEAR_SSH_PASSWORD"); len(env_gear_ssh_password) != 0 {
			password = env_gear_ssh_password
		}
		s.Config.Auth = []ssh.AuthMethod{ssh.Password(password)}
	}
	// ignore if host not in known_hosts?
	s.Config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	// SSH hostname
	var hostname string
	if gear_ssh_hostname := s.Host.Getvar("gear_ssh_hostname"); gear_ssh_hostname != nil {
		hostname = fmt.Sprint(gear_ssh_hostname)
	} else if env_gear_ssh_hostname := os.Getenv("GEAR_SSH_HOSTNAME"); len(env_gear_ssh_hostname) != 0 {
		hostname = env_gear_ssh_hostname
	} else {
		hostname = s.Host.Name
	}
	// SSH port
	var port string
	if gear_ssh_port := s.Host.Getvar("gear_ssh_port"); gear_ssh_port != nil {
		port = fmt.Sprint(gear_ssh_port)
	} else if env_gear_ssh_port := os.Getenv("GEAR_SSH_PORT"); len(env_gear_ssh_port) != 0 {
		port = env_gear_ssh_port
	} else {
		port = "22"
	}
	// Attempt connection
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), s.Config)
	if err != nil {
		return fmt.Errorf("failed to connect to host %s: %v", s.Host.Name, err)
	}
	// Create session, set environment and return session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session for host %s: %v", s.Host.Name, err)
	}
	s.Session = session
	return nil
}

func (s *SshConnection) Close() error {
	if s.Client != nil {
		err := s.Client.Close()
		if err != nil {
			return err
		}
	}
	if s.Session != nil {
		err := s.Session.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SshConnection) WhoAmI() (string, error) {
	user := s.Config.User
	if len(user) == 0 {
		return user, fmt.Errorf("could not get host username")
	}
	return user, nil
}

func (s *SshConnection) Execute(command string) *exchange.TaskResponse {
	const line_break = "-------------------------"
	getFilteredOut := func(output string) string {
		index := strings.Index(output, line_break)
		if index != -1 {
			return output[index+len(line_break):]
		}
		return output
	}
	responseErr := func(response *exchange.TaskResponse, err error) *exchange.TaskResponse {
		response.Type = exchange.Failed
		response.CommandResult.Cmd = command
		response.CommandResult.Err = err.Error()
		response.CommandResult.Rc = default_exit_code
		return response
	}
	var inpipe io.WriteCloser
	var outbuf, errbuf strings.Builder
	var stdout, stderr string
	var exitcode int
	var err error

	response := exchange.NewTaskResponse()
	inpipe, err = s.Session.StdinPipe()
	if err != nil {
		return responseErr(response, err)
	}
	s.Session.Stdout = &outbuf
	s.Session.Stderr = &errbuf
	err = s.Session.Shell()
	if err != nil {
		return responseErr(response, err)
	}
	_, err = inpipe.Write([]byte(fmt.Sprintf("echo %s\n", line_break)))
	for k, v := range s.Host.Environment {
		_, err = inpipe.Write([]byte(fmt.Sprintf("export %s=%s\n", k, v)))
	}
	_, err = inpipe.Write([]byte(fmt.Sprintf(`/usr/bin/env sh -c "LANG=C %s"`, command) + "\n"))
	if err != nil {
		response.Type = exchange.Failed
		if exit_error, ok := err.(*ssh.ExitError); ok {
			exitcode = exit_error.ExitStatus()
		} else {
			exitcode = default_exit_code
		}
	} else {
		exitcode = 0
	}

	inpipe.Close()
	err = s.Session.Wait()
	if err != nil {
		return responseErr(response, err)
	}
	stdout = getFilteredOut(outbuf.String())
	stderr = errbuf.String()

	response.CommandResult = &exchange.CommandResult{
		Cmd: command,
		Out: stdout,
		Err: stderr,
		Rc:  exitcode,
	}
	return response
}

func (s *SshConnection) CopyFile(src string, dst string) error {
	sftp, err := sftp.NewClient(s.Client)
	if err != nil {
		return err
	}
	defer sftp.Close()

	srcFile, err := utils.OpenFile(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := sftp.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		return err
	}
	return nil
}

func (s *SshConnection) WriteData(data string, path string) error {
	var file *sftp.File
	sftp, err := sftp.NewClient(s.Client)
	if err != nil {
		log.Fatalf("unable to start SFTP session: %v", err)
	}
	defer sftp.Close()
	file, err = sftp.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}
