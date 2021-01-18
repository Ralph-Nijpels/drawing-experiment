package main

import (
	"log"
	"math"
	"time"

	"./model"
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

// Draw a line diagram of a single triangular Mesh
func drawMesh(mesh model.Mesh, projection matrix.Matrix, pixels []byte) {
	drawLine(projection.Mulv(mesh.GetVertex(0)), projection.Mulv(mesh.GetVertex(1)), pixels)
	drawLine(projection.Mulv(mesh.GetVertex(1)), projection.Mulv(mesh.GetVertex(2)), pixels)
	drawLine(projection.Mulv(mesh.GetVertex(2)), projection.Mulv(mesh.GetVertex(0)), pixels)
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

	// Let's define a simple box arround the origin
	box := model.NewBox(100.0, 50.0, 25.0)
	angle := 0.0

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

		// rotate the box
		rotate := matrix.FilledMatrix([][]float32{
			{float32(math.Cos(angle)), float32(math.Sin(angle)), 0.0},
			{float32(-math.Sin(angle)), float32(math.Cos(angle)), 0.0},
			{0.0, 0.0, 1.0},
		})
		box.SetRotation(rotate)

		// move the box
		move := vector.FilledVector([]float32{
			float32(winWidth / 2), float32(0.0), float32(winHeight / 2),
		})
		box.SetPosition(move)

		// Create projection matrix
		sin30 := float32(0.5)
		cos30 := float32(0.866025)
		project := matrix.FilledMatrix([][]float32{
			{1.0, cos30, 0.0},
			{0.0, sin30, 1.0},
		})

		// Draw box 2.0
		for index := 0; index < box.Meshes(); index++ {
			drawMesh(box.GetMesh(index), project, pixels)
		}

		// Update rotation 2deg per frame
		angle += 2.0 * (math.Pi * 2.0 / 360.0)
		if angle > 2.0*math.Pi {
			angle -= 2.0 * math.Pi
		}

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
