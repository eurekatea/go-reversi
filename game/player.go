package game

import (
	"fmt"
	"os"
	"os/exec"
	"othello/board"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type player interface {
	move()
	isDone() bool
}

type human struct {
	bd    *board.Board
	color board.Color
	done  bool
}

func newHuman(bd *board.Board, cl board.Color) *human {
	return &human{bd: bd, color: cl, done: false}
}

func (h *human) move() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		x = int(float64(x-MARGIN_X)/SPACE + FIX)
		y = int(float64(y-MARGIN_Y)/SPACE + FIX)

		p := board.NewPoint(x, y)
		if h.bd.PutPoint(h.color, p) {
			h.done = true
		}
	}
}

func (h *human) isDone() bool {
	if h.done {
		h.done = false
		return true
	} else {
		return false
	}
}

type com struct {
	bd      *board.Board
	color   board.Color
	result  chan string
	ran     bool
	id      string
	program string
}

func newCom(bd *board.Board, cl board.Color, name string) *com {
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

func (c *com) isDone() bool {
	select {
	case output := <-c.result:
		col, row := (output[0] - 'A'), (output[1] - 'a')
		p := board.Point{X: int(row), Y: int(col)}
		if !c.bd.PutPoint(c.color, p) {
			r := fmt.Sprintf("this place <%s> was not valid\n", output[:2])
			r += c.bd.Visualize()
			c.fatal(r)
		}
		c.ran = false
		return true
	default:
		return false
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
	if len(output) != 3 && len(output) != 2 {
		c.fatal("unknown output: " + output)
	}

	c.result <- output
}

func (c com) fatal(text string) {
	f, err := os.Create("error.log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	text = time.Now().String() + "\n" + text

	_, err = f.Write([]byte(text))
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
