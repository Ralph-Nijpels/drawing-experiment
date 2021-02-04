package render

import (
	"fmt"
	"log"

	"../model"
	"../number/vector"
)

// Color ...
// Replace with standard GOLANG Color please
// Going to be part of the Model package, and will include many more properties
type Color struct {
	r, g, b byte
}

// Canvas defines a drawing area compatible with SDL 2.0
type Canvas struct {
	width  int
	height int
	pixels []byte
}

// NewCanvas creates a Canvas
func NewCanvas(width int, height int) *Canvas {
	return &Canvas{width, height, make([]byte, width*height*4)}
}

// Clear ...
func (c *Canvas) Clear() {
	for i := range c.pixels {
		c.pixels[i] = 0
	}
}

// Set ...
func (c *Canvas) Set(x int, y int, color Color) {
	index := (y*c.width + x) * 4
	if index < len(c.pixels)-4 && index >= 0 {
		c.pixels[index] = color.r
		c.pixels[index+1] = color.g
		c.pixels[index+2] = color.b
	}
}

// Width ...
func (c *Canvas) Width() int {
	return c.width
}

// Height ...
func (c *Canvas) Height() int {
	return c.height
}

// Pixels ...
func (c *Canvas) Pixels() []byte {
	return c.pixels
}

// Camera ...
type Camera struct {
	position vector.Vector
	lookat   vector.Vector // TODO: should be normalized
}

// NewCamera creates a camera
func NewCamera(postion vector.Vector, lookat vector.Vector) *Camera {
	return &Camera{postion, lookat}
}

// Project translates a point into the view of the camera
// The new vector represents values between ([0..1], [0..1], distance)
func (c *Camera) Project(point vector.Vector) vector.Vector {
	// Translate the Camera Direction relative to the camera position
	// giving a vector perpendicular the the eventual camera plain
	if c.position.Equal(c.lookat) {
		log.Fatalf("Render.Project: Camera position is the same as camera lookat")
	}
	cNormal := c.lookat.Sub(c.position).Unit()
	fmt.Printf("Camera normal: |%v| = %f\n", cNormal, cNormal.Abs())

	// Define the camera plane by creating 3 vectors:
	// 1 a Vector that defines the unit vector in the y-direction on this
	//   plane by rotating the normal 90deg

	// 1 a Vector that definines the position of the plane as 1 unit behind
	//   the current camera postion.
	// cSupport := c.position.Sub(cNormal)
	// Normal rotated arround the y-axis to cancel out x

	// 2 a Vector on the xy-plane that defines the x direction coordinate
	//   for any projected point
	xProjection := vector.NewVector([]float32{cNormal.Get(0).(float32), cNormal.Get(1).(float32), 0.0}).Unit()
	xCoordinate := -xProjection.Get(1).(float32) * xProjection.Get(1).(float32)
	yCoordinate := xProjection.Get(1).(float32) * xProjection.Get(0).(float32)
	xUnit := vector.NewVector([]float32{xCoordinate, yCoordinate, 0.0}).Unit()
	fmt.Printf("Unit vector for X: |%v| = %f\n", xUnit, xUnit.Abs())

	// 3 the Normal of the plane defined by the origin, the camera normal and the 
	//   xUnit vector defines the yUnit vector for the view plane
	
	// Translate the line of sight from the point to a vector representation
	// take care of the edge case where the point happens to be on the camera
	if point.Equal(c.position) {
		return vector.ZeroVector(point.Len(), point.Kind())
	}
	lnormal := c.position.Sub(point).Unit()

	// Project the viewline onto the view plane
	t := float32((cNormal.Mulv(c.position) - cNormal.Mulv(point)) / cNormal.Mulv(lnormal))
	r := point.Add(lnormal.Muls(t))

	// Translate the intersection to an x,y & depth coordinate
	return r
}

// cringeworthy version
func drawLine(from vector.Vector, to vector.Vector, canvas *Canvas, color Color) {
	// create the direction of travel and run allong the line
	dv := to.Sub(from).Unit()
	for p := from; p.Sub(to).Abs() >= 1.0; p = p.Add(dv) {
		x := int(p.Get(0).(float32) + float32(canvas.Width())/2.0)
		y := int(-p.Get(1).(float32) + float32(canvas.Height())/2.0)
		if x >= 0 && x < canvas.Width() && y >= 0 && y < canvas.Height() {
			canvas.Set(x, y, color)
		}
	}
}

// Draw a line diagram of a single triangular Mesh
func drawMesh(mesh model.Mesh, camera *Camera, canvas *Canvas) {
	drawLine(camera.Project(mesh.GetVertex(0)), camera.Project(mesh.GetVertex(1)), canvas, Color{0xff, 0xff, 0xff})
	drawLine(camera.Project(mesh.GetVertex(1)), camera.Project(mesh.GetVertex(2)), canvas, Color{0xff, 0xff, 0xff})
	drawLine(camera.Project(mesh.GetVertex(2)), camera.Project(mesh.GetVertex(0)), canvas, Color{0xff, 0xff, 0xff})
}

// For testing only
func drawGrid(camera *Camera, canvas *Canvas) {
	oo := vector.NewVector([]float32{0.0, 0.0, 0.0})
	xp := vector.NewVector([]float32{500.0, 0.0, 0.0})
	xn := vector.NewVector([]float32{-500.0, 0.0, 0.0})
	yp := vector.NewVector([]float32{0.0, 500.0, 0.0})
	yn := vector.NewVector([]float32{0.0, -500.0, 0.0})
	zp := vector.NewVector([]float32{0.0, 0.0, 500.0})
	zn := vector.NewVector([]float32{0.0, 0.0, -500.0})

	drawLine(camera.Project(oo), camera.Project(xp), canvas, Color{0xff, 0x00, 0x00})
	drawLine(camera.Project(oo), camera.Project(xn), canvas, Color{0x80, 0x80, 0x80})
	drawLine(camera.Project(oo), camera.Project(yp), canvas, Color{0x00, 0xff, 0x00})
	drawLine(camera.Project(oo), camera.Project(yn), canvas, Color{0x80, 0x80, 0x80})
	drawLine(camera.Project(oo), camera.Project(zp), canvas, Color{0x00, 0x00, 0xff})
	drawLine(camera.Project(oo), camera.Project(zn), canvas, Color{0x80, 0x80, 0x80})
}

// Draw puts the model on the canvas given current Camera and Lighting settings
func Draw(meshes []model.Mesh, camera *Camera, canvas *Canvas) {

	canvas.Clear()

	// Run trough the meshes
	drawGrid(camera, canvas)
	for _, mesh := range meshes {
		drawMesh(mesh, camera, canvas)
	}
}
