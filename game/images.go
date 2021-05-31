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
	current    *ebiten.Image
	over       *ebiten.Image
	blackWon   *ebiten.Image
	whiteWon   *ebiten.Image
	gameDraw   *ebiten.Image
)

func init() {
	var temp image.Image
	_, backGround = imageFromFS(BOARDPIC)
	temp, blackImg = imageFromFS("img/black.webp")
	_, whiteImg = imageFromFS("img/white.webp")
	_, possible = imageFromFS("img/possible.webp")
	_, current = imageFromFS("img/current.webp")
	_, over = imageFromFS(GAMEOVERPIC)
	_, blackWon = imageFromFS(BLACKWONPIC)
	_, whiteWon = imageFromFS(WHITEWONPIC)
	_, gameDraw = imageFromFS(DRAWPIC)
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

func (g *game) drawImageWithPos(screen *ebiten.Image, i, j int, draw *ebiten.Image) {
	x := float64(i)*SPACE + MARGIN_X
	y := float64(j)*SPACE + MARGIN_Y

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(x, y)

	screen.DrawImage(draw, opts)
}

func (g *game) drawStones(screen *ebiten.Image) {
	for i := 0; i < board.BOARD_LEN; i++ {
		for j := 0; j < board.BOARD_LEN; j++ {
			cl := g.bd.AtXY(i, j)
			if cl == board.BLACK {
				g.drawImageWithPos(screen, i, j, blackImg)
			} else if cl == board.WHITE {
				g.drawImageWithPos(screen, i, j, whiteImg)
			}
		}
	}
	for _, v := range g.available {
		g.drawImageWithPos(screen, v.X, v.Y, possible)
	}
	g.drawImageWithPos(screen, g.lastMove.X, g.lastMove.Y, current)
}
