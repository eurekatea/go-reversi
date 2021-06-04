package game

import (
	"fmt"
	"os"
	"os/exec"
	"othello/board"
	"strings"
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
	output, err := c.execute(bd)
	if err != nil {
		return board.Point{}, err
	}
	col, row := int(output[0]-'A'), int(output[1]-'a')
	p := board.NewPoint(row, col)
	if !bd.Put(c.color, p) {
		r := fmt.Sprintf("this place <%s> was not valid\n", output[:2])
		return board.Point{}, c.fatal(bd, r)
	}
	return p, nil
}

func (c com) execute(bd board.Board) (string, error) {
	cmd := exec.Command(c.program, "")
	cmd = modifyCmd(cmd)
	cmd.Stdin = strings.NewReader(bd.String() + c.id)
	out, err := cmd.Output()
	if err != nil {
		return "", c.fatal(bd, err.Error())
	}

	output := string(out)
	if len(output) == 0 {
		return "", c.fatal(bd, "unknown output: (no output)")
	}
	if c.invalid(bd, output) {
		return "", c.fatal(bd, "unknown output: "+output)
	}

	return output, nil
}

func (c com) invalid(bd board.Board, output string) bool {
	l := len(output) < 2
	if l {
		return true
	}
	first := output[0] < 'A' || output[0] > byte('A'+bd.Size())
	second := output[1] < 'a' || output[1] > byte('a'+bd.Size())
	return first || second
}

func (c com) fatal(bd board.Board, text string) error {
	f, err := os.Create("error.log")
	if err != nil {
		return err
	}
	defer f.Close()

	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	t := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
	text = t + "\n" + c.program + "\n\n" + text
	if len(text) > 500 {
		text = text[:500]
		text += "\n...skipped"
	}
	text += "\n\n"

	text += "last state of board:\n"
	text += bd.Visualize() + "\n"
	text += "last stdin:\n"
	text += bd.String() + c.id

	_, err = f.Write([]byte(text))
	if err != nil {
		return err
	}

	return fmt.Errorf("external AI has occured an error\nplease check the log file\nprogram exit now")
}
