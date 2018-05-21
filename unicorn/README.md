A go library for interacting with Raspberry Pi Unicorn HAT HD, with a cimilar API to the official library https://github.com/pimoroni/unicorn-hat-hd

Includes a simulator when running on amd64 and 386 linux PCs.

Sample usage:
```
package main

import "github.com/MFAshby/unicornpaint/unicorn"

func main() {
	uni, _ := unicorn.GetUnicorn()
	defer func() { uni.Off() }()
	var x, r, y, g, b uint8
	x = 10
	y = 10
	r = 255
	g = 255
	b = 255
	uni.SetPixel(x, y, r, g, b)
	uni.Show()
	uni.MainLoop()
}
```