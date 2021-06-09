package game

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"othello/board"
	"time"
)

type computer interface {
	Move(string) (string, error)
	Close()
}

type com struct {
	color board.Color
	id    string
	cmd   *exec.Cmd
	in    io.WriteCloser
	out   io.ReadCloser
}

func newCom(cl board.Color, name string) *com {
	var err error
	c := &com{
		color: cl,
		cmd:   exec.Command(name, ""),
	}
	c.cmd = modifyCmd(c.cmd)
	c.in, err = c.cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	c.out, err = c.cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	c.cmd.Start()
	if cl == board.BLACK {
		c.id = " 1"
	} else {
		c.id = " 2"
	}
	return c
}

func (c *com) Move(input string) (string, error) {
	output, err := c.execute(input)
	if err != nil {
		return "", err
	}

	bd := board.NewBoardFromStr(input)
	if !bd.PutStr(c.color, output) {
		r := fmt.Sprintf("this place <%s> was not valid\n", output[:2])
		return "", c.fatal(input, r)
	}
	return output, nil
}

func (c com) execute(input string) (string, error) {

	c.in.Write([]byte(input + c.id + "\n"))
	r := bufio.NewReader(c.out)

	var output string
	x, _ := r.ReadString('\n')
	if x != "" {
		output = string(x)
	}

	if len(output) == 0 {
		return "", c.fatal(input, "unknown output: (no output)")
	}
	if c.invalid(input, output) {
		return "", c.fatal(input, "unknown output: "+output)
	}

	return output, nil
}

func (c com) invalid(input string, output string) bool {
	if len(output) < 2 {
		return true
	}
	size := 8
	if len(input) == 6*6 {
		size = 6
	}
	first := output[0] < 'A' || output[0] > byte('A'+size)
	second := output[1] < 'a' || output[1] > byte('a'+size)
	return first || second
}

func (c com) fatal(input string, text string) error {
	bd := board.NewBoardFromStr(input)

	f, err := os.Create("error.log")
	if err != nil {
		return err
	}
	defer f.Close()

	now := time.Now()
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	t := fmt.Sprintf("%d/%02d/%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
	text = t + "\n" + c.cmd.Path + "\n\n" + text
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
