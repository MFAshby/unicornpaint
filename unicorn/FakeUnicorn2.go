package unicorn

import (
	"image"

	"github.com/veandco/go-sdl2/sdl"
)

type FakeUnicorn2 struct {
	BaseUnicorn2
}

func renderImage(im image.Image) {
	// b := im.Bounds()
	// width := b.Dx()
	// height := b.Dy()
	// for x := 0; x < width; x++ {
	// 	for y := 0; y < height; y++ {
	// 		r, g, b, _ := im.At(x, y).RGBA()
	// 		un.SetPixel(uint8(x), uint8(y), uint8(r), uint8(g), uint8(b))
	// 	}
	// }
	// un.Show()
}

func render() {
	// for !stop {
	// 	for i := 0; i < len(gf.Image); i++ {
	// 		im := gf.Image[i]
	// 		delay := gf.Delay[i] //100ths of a second
	// 		renderImage(un, im)
	// 		time.Sleep(time.Duration(delay * 10000000)) // nanoseconds 10^-9 sec
	// 	}
	// }
}

func (u *FakeUnicorn2) StartRender() {

}

func MainLoop() {
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func NewUnicorn2() *FakeUnicorn2 {
	return &FakeUnicorn2{}
}
