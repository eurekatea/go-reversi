//go:build !windows

package game

import (
	"os/exec"
	"syscall"
)

func modifyCmd(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Setpgid = true
	return cmd
}

func (c com) Close() {
	err := syscall.Kill(-c.cmd.Process.Pid, syscall.SIGKILL)
	if err != nil {
		panic(err)
	}
}
