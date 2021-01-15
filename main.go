package main

import (
	"log"
	"time"

	"./number/matrix"
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
	for p := from; p.Sub(to).Abs() >= 1.0; p = p.Add(dv) {
		setPixel(int(p.Get(0).(float32)), int(p.Get(1).(float32)), color{0xFF, 0xFF, 0xFF}, pixels)
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
	ftl := vector.FilledVector([]float32{100.0, 100.0, 100.0})
	ftr := vector.FilledVector([]float32{300.0, 100.0, 100.0})
	fbr := vector.FilledVector([]float32{300.0, 300.0, 100.0})
	fbl := vector.FilledVector([]float32{100.0, 300.0, 100.0})
	btl := vector.FilledVector([]float32{100.0, 100.0, 300.0})
	btr := vector.FilledVector([]float32{300.0, 100.0, 300.0})
	bbr := vector.FilledVector([]float32{300.0, 300.0, 300.0})
	bbl := vector.FilledVector([]float32{100.0, 300.0, 300.0})

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
		sin30 := float32(0.5)
		cos30 := float32(0.866025)
		project := matrix.FilledMatrix([][]float32{
			{1.0, cos30, 0.0},
			{0.0, sin30, 1.0},
		})

		// Draw box
		drawLine(project.Mulv(ftl), project.Mulv(ftr), pixels)
		drawLine(project.Mulv(ftr), project.Mulv(fbr), pixels)
		drawLine(project.Mulv(fbr), project.Mulv(fbl), pixels)
		drawLine(project.Mulv(fbl), project.Mulv(ftl), pixels)

		drawLine(project.Mulv(ftl), project.Mulv(btl), pixels)
		drawLine(project.Mulv(ftr), project.Mulv(btr), pixels)
		drawLine(project.Mulv(fbr), project.Mulv(bbr), pixels)
		drawLine(project.Mulv(fbl), project.Mulv(bbl), pixels)

		drawLine(project.Mulv(btl), project.Mulv(btr), pixels)
		drawLine(project.Mulv(btr), project.Mulv(bbr), pixels)
		drawLine(project.Mulv(bbr), project.Mulv(bbl), pixels)
		drawLine(project.Mulv(bbl), project.Mulv(btl), pixels)

		// Rotate box
		sin2 := float32(0.0348999)
		cos2 := float32(0.9999391)
		rotate := matrix.FilledMatrix([][]float32{
			{cos2, sin2, 0.0},
			{-sin2, cos2, 0.0},
			{0.0, 0.0, 1.0},
		})

		ftl = rotate.Mulv(ftl)
		ftr = rotate.Mulv(ftr)
		fbr = rotate.Mulv(fbr)
		fbl = rotate.Mulv(fbl)
		btl = rotate.Mulv(btl)
		btr = rotate.Mulv(btr)
		bbr = rotate.Mulv(bbr)
		bbl = rotate.Mulv(bbl)

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
