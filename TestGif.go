// +build ignore

package main

import (
	"bytes"
	"image/gif"
	"log"

	"github.com/MFAshby/unicornpaint/unicorn"
)

func main() {
	un, err := unicorn.NewUnicorn2()
	if err != nil {
		log.Fatalf("Error getting a unicorn :( %v", err)
	}

	data, err := unicorn.Asset("data/sample2.gif")
	if err != nil {
		log.Fatalf("Error getting rain %v", err)
	}

	g, err := gif.DecodeAll(bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Error decoding gif %v", err)
	}

	un.SetGif(g)
	stopChan := un.StartRender()
	un.MainLoop()
	stopChan <- true
}
