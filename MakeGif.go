package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io/ioutil"
	"time"

	"github.com/MFAshby/unicornpaint/unicorn"
)

var (
	un unicorn.Unicorn
)

func imageFromPixels(pixels [][][]uint8) image.Image {
	width := len(pixels)
	height := len(pixels[0])
	im1 := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b := unicorn.Rgb(pixels[x][y])
			col := color.RGBA{
				R: r,
				G: g,
				B: b,
				A: 100,
			}
			im1.Set(x, y, col)
		}
	}
	return im1
}

func toPaletted(im image.Image) *image.Paletted {
	b := im.Bounds()
	pm := image.NewPaletted(b, palette.Plan9[:256])
	draw.FloydSteinberg.Draw(pm, b, im, image.ZP)
	return pm
}

var (
	stop bool
)

func renderImage(un unicorn.Unicorn, im image.Image) {
	b := im.Bounds()
	width := b.Dx()
	height := b.Dy()
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, _ := im.At(x, y).RGBA()
			un.SetPixel(uint8(x), uint8(y), uint8(r), uint8(g), uint8(b))
		}
	}
	un.Show()
}

func renderGif(un unicorn.Unicorn, gf *gif.GIF) {
	for !stop {
		for i := 0; i < len(gf.Image); i++ {
			im := gf.Image[i]
			delay := gf.Delay[i] //100ths of a second
			renderImage(un, im)
			time.Sleep(time.Duration(delay * 10000000)) // nanoseconds 10^-9 sec
		}
	}
}

func main() {
	b1, _ := ioutil.ReadFile("saves/modern")
	b2, _ := ioutil.ReadFile("saves/modern2")

	px1 := [][][]uint8{}
	json.Unmarshal(b1, &px1)
	px2 := [][][]uint8{}
	json.Unmarshal(b2, &px2)

	im1 := toPaletted(imageFromPixels(px1))
	im2 := toPaletted(imageFromPixels(px2))

	gf := &gif.GIF{
		Image: []*image.Paletted{im1, im2},
		Delay: []int{50, 50}, // 100ths of a second
	}

	// f1, err := os.Create("saves/modern.gif")
	// if err != nil {
	// 	log.Fatalf("Error opening GIF file to write %v", err)
	// }
	// defer f1.Close()
	// err = gif.EncodeAll(f1, gf)
	// if err != nil {
	// 	log.Printf("Error writing GIF %v", err)
	// }

	un, _ = unicorn.NewUnicorn()

	go renderGif(un, gf)

	un.MainLoop()
}
