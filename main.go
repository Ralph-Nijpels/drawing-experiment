package main

import (
	"log"
	"time"

	"./number/vector"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 800
	winHeight = 600
)

// Replace with standard GOLANG Color please
type color struct {
	r, g, b byte
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

// cringeworthy version
func drawLine(from vector.Vector, to vector.Vector, pixels []byte) {
	// create the direction of travel and run allong the line
	dv := to.Sub(from).Unit()
	for p := from; p.Sub(to).Abs().(float64) >= 1.0; p = p.Add(dv) {
		setPixel(int(p.Get(0).(float64)), int(p.Get(1).(float64)), color{0xFF, 0xFF, 0xFF}, pixels)
	}
}

func main() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Panic(err)
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		log.Panic(err)
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	// Initialize players

	// Get initial key states
	// keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32

	// Not so nice wat to establish points, we have to update 'vector' a litle. To bad I did not understand a thing about
	// reflection at that point
	topLeft := vector.Make([]float64{100.0, 100.0})
	topRight := vector.Make([]float64{400.0, 100.0})
	bottomRight := vector.Make([]float64{400.0, 200.0})
	bottomLeft := vector.Make([]float64{100.0, 200.0})

	// Big game loop
	for {
		frameStart = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)

		// Update elements
		drawLine(topLeft, topRight, pixels)
		drawLine(topRight, bottomRight, pixels)
		drawLine(bottomRight, bottomLeft, pixels)
		drawLine(bottomLeft, topLeft, pixels)
		drawLine(topLeft, bottomRight, pixels)
		drawLine(bottomLeft, topRight, pixels)

		// Show results
		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		// Keep a steady flow
		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < 0.005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}

}
