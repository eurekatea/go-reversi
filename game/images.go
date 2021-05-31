package game

import (
	"embed"
	"image"
	"othello/board"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/webp"
)

//go:embed img
var source embed.FS

var (
	icon       []image.Image
	backGround *ebiten.Image
	blackImg   *ebiten.Image
	whiteImg   *ebiten.Image
	possible   *ebiten.Image
	over       *ebiten.Image
	blackWon   *ebiten.Image
	whiteWon   *ebiten.Image
	gameDraw   *ebiten.Image
)

func init() {
	var temp image.Image
	_, backGround = imageFromFS("img/board6x6.webp")
	temp, blackImg = imageFromFS("img/black.webp")
	_, whiteImg = imageFromFS("img/white.webp")
	_, possible = imageFromFS("img/possible.webp")
	_, over = imageFromFS("img/gameover.webp")
	_, blackWon = imageFromFS("img/blackwon.webp")
	_, whiteWon = imageFromFS("img/whitewon.webp")
	_, gameDraw = imageFromFS("img/gamedraw.webp")
	icon = []image.Image{temp}
}

func Icon() []image.Image {
	return icon
}

func imageFromFS(path string) (image.Image, *ebiten.Image) {
	f, err := source.Open(path)
	if err != nil {
		panic(err)
	}
	bytes, err := webp.Decode(f)
	if err != nil {
		panic(err)
	}
	img := ebiten.NewImageFromImage(bytes)
	return bytes, img
}

func (g *game) drawBoard(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	screen.DrawImage(backGround, options)
}

func (g *game) drawEnd(screen *ebiten.Image) {
	options := &ebiten.DrawImageOptions{}
	screen.DrawImage(over, options)

	if g.winner == board.BLACK {
		screen.DrawImage(blackWon, options)
	} else if g.winner == board.WHITE {
		screen.DrawImage(whiteWon, options)
	} else {
		screen.DrawImage(gameDraw, options)
	}
}

const (
	SPACE    = 59  // the SPACE between every stone
	MARGIN_X = 42  // for the first stone
	MARGIN_Y = 42  // for the first stone
	FIX      = 0.1 // FIX the position inaccuracy
)

func (g *game) drawImageFromPoint(screen *ebiten.Image, p board.Point, draw *ebiten.Image) {
	x := float64(p.X)*SPACE + MARGIN_X
	y := float64(p.Y)*SPACE + MARGIN_Y

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)

	screen.DrawImage(draw, opts)
}

func (g *game) drawStones(screen *ebiten.Image) {
	for i := 0; i < board.BOARD_LEN; i++ {
		for j := 0; j < board.BOARD_LEN; j++ {
			p := board.NewPoint(i, j)
			cl := g.bd.AtPoint(p)
			if cl == board.NONE {
				continue
			} else if cl == board.BLACK {
				g.drawImageFromPoint(screen, p, blackImg)
			} else if cl == board.WHITE {
				g.drawImageFromPoint(screen, p, whiteImg)
			}
		}
	}

	for _, v := range g.available {
		g.drawImageFromPoint(screen, v, possible)
	}
}
