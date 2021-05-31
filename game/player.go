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
	g     *game
	bd    *board.Board
	color board.Color
	done  bool
}

func newHuman(g *game, bd *board.Board, cl board.Color) *human {
	return &human{g: g, bd: bd, color: cl, done: false}
}

func (h *human) move() {
	x, y := ebiten.CursorPosition()

	x = int(float64(x-MARGIN_X)/SPACE + FIX)
	y = int(float64(y-MARGIN_Y)/SPACE + FIX)

	h.hint(x, y)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p := board.NewPoint(x, y)
		if h.bd.Put(h.color, p) {
			h.done = true
			ebiten.SetCursorShape(ebiten.CursorShapeDefault)
		}
	}
}

func (h *human) hint(x, y int) {
	var need = false
	for _, p := range h.g.available {
		if p.X == x && p.Y == y {
			need = true
		}
	}
	if need {
		ebiten.SetCursorShape(ebiten.CursorShapePointer)
	} else {
		ebiten.SetCursorShape(ebiten.CursorShapeDefault)
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
		col, row := int(output[0]-'A'), int(output[1]-'a')
		p := board.NewPoint(row, col)
		if !c.bd.Put(c.color, p) {
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
	if c.invalid(output) {
		c.fatal("unknown output: " + output)
	}

	c.result <- output
}

func (c com) invalid(output string) bool {
	l := len(output) < 2
	first := output[0] < 'A' || output[0] > 'A'+board.BOARD_LEN
	second := output[1] < 'a' || output[1] > 'a'+board.BOARD_LEN
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
		text += "\n...skipped\n"
	}

	_, err = f.Write([]byte(text))
	if err != nil {
		panic(err)
	}

	panic("error")
}
