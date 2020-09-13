package backend

import (
	"strings"

	"github.com/raba-jp/primus/pkg/internal/exec"
	"github.com/spf13/afero"
)

type OS int

const (
	Unknown OS = iota + 1
	Darwin
	Arch
)

func DetectOS(execIF exec.Interface, fs afero.Fs) OS {
	if darwin := DetectDarwin(execIF); darwin {
		return Darwin
	}
	if manjaro := DetectArchLinux(fs); manjaro {
		return Arch
	}
	return Unknown
}

func DetectDarwin(execIF exec.Interface) bool {
	out, err := execIF.Command("uname", "-a").Output()
	if err != nil {
		return false
	}
	if !strings.Contains(string(out), "Darwin") {
		return false
	}
	return true
}

func DetectArchLinux(fs afero.Fs) bool {
	_, err := fs.Stat("/etc/arch-release")
	return err == nil
}
