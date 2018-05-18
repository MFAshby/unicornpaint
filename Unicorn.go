package main

import (
	"github.com/veandco/go-sdl2/sdl"
	// "golang.org/x/exp/io/spi"
	// "github.com/veandco/go-sdl2/sdl"
)

// Unicorn ...
// Object representing the Unicorn HAT to be controlled
type Unicorn interface {
	// Not all unicorns are the same size
	GetWidth() uint8
	GetHeight() uint8

	// Array of pixels, indexed x, then y, then color (rgb)
	GetPixels() [][][]uint8

	// Set an individual pixel
	SetPixel(x, y, r, g, b uint8)

	// Flip the display buffer
	Show()

	// Set all pixels back to black
	Clear()

	// Turns off the LEDs
	Off()
}

// GetUnicorn ...
// Get a unicorn. Tries to get you a real one,
// if it can't find one then gives you a fake one.
func GetUnicorn() (unicorn Unicorn, err error) {
	// unicorn, err = NewReal()
	// if err != nil {
	// 	log.Println("Couldn't get a real unicorn, trying a fake one")
	// 	unicorn, err = NewFake(int8(16), int8(16))
	// }
	unicorn, err = NewFake(uint8(16), uint8(16))
	return
}

// FakeUnicorn ...
// Shows an SDL window pretending to be a unicorn.
type FakeUnicorn struct {
	pixels        [][][]uint8
	displayWidth  int32
	displayHeight int32

	window   *sdl.Window
	renderer *sdl.Renderer
}

// NewFake ...
// Constructs a new fake unicorn out of paint and glue
func NewFake(width, height uint8) (*FakeUnicorn, error) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return nil, err
	}

	unicorn := &FakeUnicorn{
		pixels:        makePixels(width, height),
		window:        nil,
		renderer:      nil,
		displayWidth:  300,
		displayHeight: 300,
	}
	if err := unicorn.createWindow(); err != nil {
		unicorn.Close()
		return nil, err
	}
	if err := unicorn.createRenderer(); err != nil {
		unicorn.Close()
		return nil, err
	}
	return unicorn, nil
}

func (f *FakeUnicorn) createWindow() error {
	window, err := sdl.CreateWindow("Fake Unicorn",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		f.displayWidth,
		f.displayHeight,
		sdl.WINDOW_SHOWN)
	f.window = window
	return err
}

func (f *FakeUnicorn) createRenderer() error {
	renderer, err := sdl.CreateRenderer(f.window, -1, sdl.RENDERER_ACCELERATED)
	f.renderer = renderer
	return err
}

func (f *FakeUnicorn) Close() error {
	if f.window != nil {
		f.window.Destroy()
	}
	if f.renderer != nil {
		f.renderer.Destroy()
	}
	return nil
}

func (f *FakeUnicorn) GetWidth() uint8 {
	return uint8(len(f.pixels))
}
func (f *FakeUnicorn) GetHeight() uint8 {
	if len(f.pixels) > 0 {
		return uint8(len(f.pixels[0]))
	}
	return 0
}
func (f *FakeUnicorn) GetPixels() [][][]uint8 {
	return f.pixels
}
func (f *FakeUnicorn) SetPixel(x, y, r, g, b uint8) {
	f.pixels[x][y] = []uint8{r, g, b}
}
func (f *FakeUnicorn) Show() {
	width, height := f.GetWidth(), f.GetHeight()
	for x := uint8(0); x < width; x++ {
		for y := uint8(0); y < height; y++ {
			r, g, b := rgb(f.pixels[x][y])
			if err := f.renderer.SetDrawColor(r, g, b, uint8(255)); err != nil {
				panic(err)
			}
			cellWidth := f.displayWidth / int32(width)
			cellHeight := f.displayHeight / int32(height)
			if err := f.renderer.FillRect(&sdl.Rect{
				X: cellWidth * int32(x),
				Y: f.displayHeight - (cellHeight * int32(y)) - cellHeight, // SDL Y coordinate is from the top
				W: cellWidth,
				H: cellHeight,
			}); err != nil {
				panic(err)
			}
		}
	}
	f.renderer.Present()
}

func rgb(pixel []uint8) (uint8, uint8, uint8) {
	return pixel[0], pixel[1], pixel[2]
}

func (f *FakeUnicorn) Clear() {
	f.pixels = makePixels(f.GetWidth(), f.GetHeight())
}
func (f *FakeUnicorn) Off() {
	f.Close()
}

func makePixels(width, height uint8) [][][]uint8 {
	pixels := make([][][]uint8, width)
	for x := uint8(0); x < width; x++ {
		pixels[x] = make([][]uint8, height)
		for y := uint8(0); y < height; y++ {
			pixels[x][y] = []uint8{0, 0, 0}
		}
	}
	return pixels
}

// RealUnicorn ...
// A real one! *gasps*
// type RealUnicorn struct {}

// // NewReal ...
// // Constructs a new real unicorn from fairy dust and sprinkles
// func NewReal() (*RealUnicorn, error) {
// 	return nil, errors.New("Couldn't make a real unicorn sorry")
// }

// func (u *RealUnicorn) GetWidth() int8 {
// 	return 0
// }
// func (u *RealUnicorn) GetHeight() int8 {
// 	return 0
// }
// func (u *RealUnicorn) GetPixels() [][][]int8 {
// 	return nil
// }
// func (u *RealUnicorn) SetPixel(x, y, r, g, b int8) {

// }
// func (u *RealUnicorn) Show() {

// }
// func (u *RealUnicorn) Clear() {

// }
// func (u *RealUnicorn) Off() {

// }
