package game

import (
	"bufio"
	"bytes"
	"embed"
	"image/png"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"golang.org/x/image/webp"
)

//go:embed img
var source embed.FS

var (
	blackImg *fyne.StaticResource
	whiteImg *fyne.StaticResource
	noneImg  *fyne.StaticResource
	possible *fyne.StaticResource
	current  *fyne.StaticResource
)

func init() {
	blackImg = resourceFromBytes("img/black.webp")
	whiteImg = resourceFromBytes("img/white.webp")
	noneImg = resourceFromBytes("img/none.webp")
	possible = resourceFromBytes("img/possible.webp")
	current = resourceFromBytes("img/current.webp")
}

func resourceFromBytes(path string) *fyne.StaticResource {
	cont := bytesFromFS(path)
	return fyne.NewStaticResource(path, cont)
}

// fyne didn't support webp so convert to png first
func bytesFromFS(path string) (cont []byte) {
	f, err := source.Open(path)
	if err != nil {
		panic(err)
	}
	img, err := webp.Decode(f)
	if err != nil {
		panic(err)
	}
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, img)
	if err != nil {
		panic(err)
	}
	cont, err = ioutil.ReadAll(bufio.NewReader(buffer))
	if err != nil {
		panic(err)
	}
	return
}
