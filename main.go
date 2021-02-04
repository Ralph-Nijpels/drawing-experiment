package main

import (
	"log"
	"time"

	"./model"
	"./number/vector"
	"./render"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	winWidth  = 800
	winHeight = 600
)

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

	// Initialize players

	// Get initial key states
	// keyState := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32

	// Let's define the rendering environment
	canvas := render.NewCanvas(winWidth, winHeight)
	camera := render.NewCamera(
		vector.NewVector([]float32{0.0, -300.0, 0.0}),
		vector.NewVector([]float32{0.0, 0.0, 0.0}))

	// Let's define a simple box arround the origin
	box := model.NewBox(100.0, 100.0, 100.0)
	angle := float32(0.0)

	// Big game loop
	for {
		frameStart = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		// rotate the box
		rotate := vector.NewVector([]float32{
			0.0, 0.0, angle,
		})
		box.SetRotation(rotate)

		// Draw box 2.0
		render.Draw(box.GetMeshes(), camera, canvas)

		// Update rotation 2deg per frame
		angle += 2.0
		if angle > 360.0 {
			angle -= 360.0
		}

		// Show results
		tex.Update(nil, canvas.Pixels(), canvas.Width()*4)
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
