package game

import (
	"fmt"
	"os"
	"os/exec"
	"othello/game/board"
	"strings"
	"time"
)

type player interface {
	move()
	isDone() (board.Point, bool)
}

type human struct {
	bd     board.Board
	color  board.Color
	done   bool
	result board.Point
}

func newHuman(bd board.Board, cl board.Color) *human {
	return &human{bd: bd, color: cl, done: false}
}

func (h *human) move() {

}

func (h *human) isDone() (board.Point, bool) {
	if h.done {
		h.done = false
		return h.result, true
	} else {
		return board.NewPoint(-1, -1), false
	}
}

type com struct {
	bd      board.Board
	color   board.Color
	result  chan string
	ran     bool
	id      string
	program string
}

func newCom(bd board.Board, cl board.Color, name string) *com {
	c := &com{
		bd:      bd,
		color:   cl,
		result:  make(chan string),
		ran:     false,
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
	if !c.ran {
		c.ran = true
		go c.execute()
	}
}

func (c *com) isDone() (board.Point, bool) {
	select {
	case output := <-c.result:
		col, row := int(output[0]-'A'), int(output[1]-'a')
		p := board.NewPoint(row, col)
		if !c.bd.Put(c.color, p) {
			r := fmt.Sprintf("this place <%s> was not valid\n", output[:2])
			r += c.bd.Visualize()
			c.fatal(r)
		}
		c.ran = false
		return p, true
	default:
		return board.NewPoint(-1, -1), false
	}
}

func (c com) execute() {
	cmd := exec.Command(execCmd+c.program, "")
	cmd.Stdin = strings.NewReader(c.bd.String() + c.id)
	out, err := cmd.Output()
	if err != nil {
		c.fatal(err.Error())
	}

	output := string(out)
	if c.invalid(output) {
		c.fatal("unknown output: " + output)
	}

	c.result <- output
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
