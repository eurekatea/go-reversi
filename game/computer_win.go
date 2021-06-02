// +build windows

package game

import (
	"fmt"
	"os"
	"os/exec"
	"othello/board"
	"strings"
	"syscall"
	"time"
)

type com struct {
	color   board.Color
	id      string
	program string
}

func newCom(cl board.Color, name string) *com {
	c := &com{
		color:   cl,
		program: name,
	}
	if cl == board.BLACK {
		c.id = " 1"
	} else {
		c.id = " 2"
	}
	return c
}

func (c *com) Move(bd board.Board) (board.Point, error) {
	output := c.execute(bd)
	col, row := int(output[0]-'A'), int(output[1]-'a')
	p := board.NewPoint(row, col)
	if !bd.Put(c.color, p) {
		r := fmt.Sprintf("this place <%s> was not valid\n", output[:2])
		r += bd.Visualize()
		return board.Point{}, c.fatal(r)
	}
	return p, nil
}

func (c com) execute(bd board.Board) string {
	cmdPath := "C:\\Windows\\system32\\cmd.exe"
	cmdInstance := exec.Command(cmdPath, c.program)
	cmdInstance.SysProcAttr = &syscall.SysProcAttr{HideWindow: false}
	cmdInstance.Stdin = strings.NewReader(bd.String() + c.id)
	out, err := cmdInstance.Output()
	if err != nil {
		c.fatal(err.Error())
	}

	output := string(out)
	if c.invalid(bd, output) {
		c.fatal("unknown output: " + output)
	}

	return output
}

func (c com) invalid(bd board.Board, output string) bool {
	l := len(output) < 2
	first := output[0] < 'A' || output[0] > byte('A'+bd.Size())
	second := output[1] < 'a' || output[1] > byte('a'+bd.Size())
	return l || first || second
}

func (c com) fatal(text string) error {
	f, err := os.Create("error.log")
	if err != nil {
		return err
	}
	defer f.Close()

	text = time.Now().String() + "\n" + text
	if len(text) > 500 {
		text = text[:500]
		text += "\n...skipped"
	}
	text += "\n"

	_, err = f.Write([]byte(text))
	if err != nil {
		return err
	}

	return fmt.Errorf("selected external AI has occured an error\nplease check the log file\nprogram will exit now")
}
