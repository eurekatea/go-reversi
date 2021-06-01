package game

import (
	"fmt"
	"os"
	"os/exec"
	"othello/game/board"
	"strings"
	"time"
)

type com struct {
	bd      board.Board
	color   board.Color
	id      string
	program string
}

func newCom(bd board.Board, cl board.Color, name string) *com {
	c := &com{
		bd:      bd,
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

func (c *com) move() {
	output := c.execute()
	col, row := int(output[0]-'A'), int(output[1]-'a')
	p := board.NewPoint(row, col)
	if !c.bd.Put(c.color, p) {
		r := fmt.Sprintf("this place <%s> was not valid\n", output[:2])
		r += c.bd.Visualize()
		c.fatal(r)
	}
}

func (c com) execute() string {
	cmd := exec.Command(c.program, "")
	cmd.Stdin = strings.NewReader(c.bd.String() + c.id)
	out, err := cmd.Output()
	if err != nil {
		c.fatal(err.Error())
	}

	output := string(out)
	if c.invalid(output) {
		c.fatal("unknown output: " + output)
	}

	return output
}

func (c com) invalid(output string) bool {
	l := len(output) < 2
	first := output[0] < 'A' || output[0] > byte('A'+c.bd.Size())
	second := output[1] < 'a' || output[1] > byte('a'+c.bd.Size())
	return l || first || second
}

func (c com) fatal(text string) {
	f, err := os.Create("error.log")
	if err != nil {
		panic(err)
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
		panic(err)
	}

	panic("error")
}
