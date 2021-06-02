// +build windows

package game

import (
	"os/exec"
	"syscall"
)

func modifyCmd(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
