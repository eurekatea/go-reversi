//go:build windows

package game

import (
	"fmt"
	"os/exec"
	"syscall"
)

func modifyCmd(cmd *exec.Cmd) *exec.Cmd {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}

func (c com) Close() {
	exec.Command("taskkill", "/F", "/T", "/PID", fmt.Sprint(c.cmd.Process.Pid)).Run()
}
