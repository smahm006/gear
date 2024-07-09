package connection

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/sftp"
	"github.com/smahm006/gear/src/inventory"
	"github.com/smahm006/gear/src/utils"
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
		Host: Host,
	}
}

func (s *SshConnection) Connect() error {
	// SSH username
	var username string
	if gsu := s.Host.Getenv("gear_ssh_user"); len(gsu) != 0 {
		username = gsu
	} else if env_gsu := os.Getenv("GEAR_SSH_USER"); len(env_gsu) != 0 {
		username = env_gsu
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
	if gspk := s.Host.Getenv("gear_ssh_private_key"); len(gspk) != 0 {
		pk_data, err := utils.ReadFile(gspk)
		if err == nil {
			signer, err = ssh.ParsePrivateKey(pk_data)
		}
	} else if env_gspk := s.Host.Getenv("GEAR_SSH_PRIVATE_KEY"); err != nil && len(env_gspk) != 0 {
		pk_data, err := utils.ReadFile(env_gspk)
		if err == nil {
			signer, err = ssh.ParsePrivateKey(pk_data)
		}
	}
	if signer != nil {
		s.Config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		// Could not find or parse a private key
		if gsp := s.Host.Getenv("gear_ssh_password"); len(gsp) != 0 {
			password = gsp
		} else if env_gsp := os.Getenv("GEAR_SSH_PASSWORD"); len(env_gsp) != 0 {
			password = env_gsp
		}
		s.Config.Auth = []ssh.AuthMethod{ssh.Password(password)}
	}
	// ignore if host not in known_hosts?
	s.Config.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	// SSH hostname
	var hostname string
	if gsh := s.Host.Getenv("gear_ssh_hostname"); len(gsh) != 0 {
		hostname = gsh
	} else if env_gsh := os.Getenv("GEAR_SSH_HOSTNAME"); len(env_gsh) != 0 {
		hostname = env_gsh
	} else {
		hostname = s.Host.Name
	}
	// SSH port
	var port string
	if gsh := s.Host.Getenv("gear_ssh_port"); len(gsh) != 0 {
		port = gsh
	} else if env_gsh := os.Getenv("GEAR_SSH_PORT"); len(env_gsh) != 0 {
		port = env_gsh
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
	for k, v := range s.Host.Environment {
		session.Setenv(k, v)
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

func (s *SshConnection) Execute(string) (string, error) {
	return "", nil
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
