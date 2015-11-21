// +build darwin freebsd linux netbsd openbsd

package term

import (
	"os"
	"os/exec"
)

func IsANSI(f *os.File) bool {
	return IsTerminal(f)
}

func IsTerminal(f *os.File) bool {
	cmd := exec.Command("test", "-t", "0")
	cmd.Stdin = f
	return cmd.Run() == nil
}
