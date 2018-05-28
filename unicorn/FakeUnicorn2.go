// +build linux,386 linux,amd64

package unicorn

import (
	"image"

	"github.com/veandco/go-sdl2/sdl"
)

type FakeUnicorn2 struct {
	BaseUnicorn2
	*BaseFakeUnicorn
}

func (u *FakeUnicorn2) renderImage(im image.Image) {
	b := im.Bounds()
	width, height := b.Dx(), b.Dy()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			col := im.At(x, y)
			r, g, b, _ := col.RGBA()
			// Ignore alpha for now, not worked out how it should work on real unicorn
			if err := u.renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), uint8(255)); err != nil {
				panic(err)
			}
			cellWidth := u.displayWidth / int32(width)
			cellHeight := u.displayHeight / int32(height)
			if err := u.renderer.FillRect(&sdl.Rect{
				X: cellWidth * int32(x),
				Y: u.displayHeight - (cellHeight * int32(y)) - cellHeight, // SDL Y coordinate is from the top
				W: cellWidth,
				H: cellHeight,
			}); err != nil {
				panic(err)
			}
		}
	}
	u.renderer.Present()
}

func (u *FakeUnicorn2) StartRender() chan bool {
	return u.StartRenderBase(u.renderImage)
}

func NewUnicorn2() (*FakeUnicorn2, error) {
	baseFake, err := NewBaseFakeUnicorn(300, 300)
	if err != nil {
		return nil, err
	}
	return &FakeUnicorn2{
		BaseUnicorn2{},
		baseFake,
	}, nil
}
