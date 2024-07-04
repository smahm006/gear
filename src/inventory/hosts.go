package inventory

import (
	"regexp"
	"strings"
)

type OS int64

const (
	Linux OS = iota
)

type PackageManager int64

const (
	Apt PackageManager = iota
	Pacman
	Rpm
)

type Host struct {
	Name           string
	Variables      map[string]interface{}
	Environment    map[string]string
	Os             *OS
	PackageManager *PackageManager
}

func NewHost(name string) *Host {
	return &Host{
		Name:           name,
		Variables:      make(map[string]interface{}),
		Environment:    make(map[string]string),
		Os:             nil,
		PackageManager: nil,
	}
}

func (h *Host) SetOs(uname string) {
	if strings.HasPrefix(uname, "Linux") {
		h.Os = new(OS)
		*h.Os = Linux
	}
}

func (h *Host) SetDistro(osrelease string) {
	apt_distros := regexp.MustCompile(`Debian|Ubuntu|LinuxMint`)
	pacman_distros := regexp.MustCompile(`Arch Linux|Manjaro`)
	rpm_distros := regexp.MustCompile(`CentOS|Fedora`)
	switch {
	case apt_distros.MatchString(osrelease):
		h.PackageManager = new(PackageManager)
		*h.PackageManager = Apt
	case pacman_distros.MatchString(osrelease):
		h.PackageManager = new(PackageManager)
		*h.PackageManager = Pacman
	case rpm_distros.MatchString(osrelease):
		h.PackageManager = new(PackageManager)
		*h.PackageManager = Rpm
	}
}

func (h *Host) IsLocal() bool {
	return h.Name == "127.0.0.1"
}

func (h *Host) Getenv(name string) string {
	env := h.Environment[name]
	return env
}
