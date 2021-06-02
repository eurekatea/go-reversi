// +build !windows

package game

import "os/exec"

func modifyCmd(cmd *exec.Cmd) *exec.Cmd {
	return cmd
}
