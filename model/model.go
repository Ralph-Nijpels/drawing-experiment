package model

import (
	"fmt"
	"log"
	"math"
	"reflect"

	"../number/matrix"
	"../number/vector"
)

// Mesh is a 3D triangle
type Mesh struct {
	vertices [3]vector.Vector
}

// NewMesh creates a Mesh
func NewMesh(vertices []vector.Vector) Mesh {
	// Check size of the polygon, there should be at least three points
	if len(vertices) != 3 {
		log.Fatalf("Model.NewMesh: expected 3 points, got %d", len(vertices))
	}

	var mesh Mesh
	for i, point := range vertices {
		if point.Len() != 3 {
			log.Fatalf("Model.NewMesh: expected 3D point, got %d", point.Len())
		}
		mesh.vertices[i] = point
	}

	return mesh
}

// GetVertex returns a vertex from a Mesh
func (m Mesh) GetVertex(index int) vector.Vector {
	if index < 0 || index > 2 {
		log.Fatalf("Model.GetVertex: index out of range %d", index)
	}
	return m.vertices[index]
}

// Add stringer interface
func (m Mesh) String() string {
	return fmt.Sprintf("[\n\t%v\n\t%v\n\t%v\n]", m.GetVertex(0), m.GetVertex(1), m.GetVertex(2))
}

// Part represents a complex object
type Part struct {
	position vector.Vector
	rotation matrix.Matrix
	scaling  matrix.Matrix
	meshes   []Mesh
}

// SetPosition moves the part arround in it's parents coordinate system
func (p *Part) SetPosition(position vector.Vector) {
	if position.Len() != 3 && position.Kind() != reflect.Float32 {
		log.Fatalf("Part.SetPosition: expects 3D-Float32 vector, got %dD-%v", position.Len(), position.Kind())
	}
	p.position = position
}

// SetRotation moves the part arround inside it's own coordinate system
func (p *Part) SetRotation(rotation vector.Vector) {
	if rotation.Len() != 3 && rotation.Kind() != reflect.Float32 {
		log.Fatalf("Part.SetRotation: expects 3D-Float32 vector, got %dD-%v", position.Len(), position.Kind())
	}
	// Make float64 values in Radians for math
	x := float64(rotation.Get(0).(float32)) * math.Pi / 180.0
	y := float64(rotation.Get(1).(float32)) * math.Pi / 180.0
	z := float64(rotation.Get(2).(float32)) * math.Pi / 180.0
	// Create the rotation matrix
	p.rotation = matrix.FilledMatrix([][]float32{
		{}
	})

}

// GetMesh returns a mesh rotated and positioned as set by the part
func (p *Part) GetMesh(index int) Mesh {
	// Check input
	if index < 0 || index > p.Meshes() {
		log.Fatalf("Part.GetMesh: index expected [0..%d>, got %d", p.Meshes(), index)
	}

	mesh := NewMesh([]vector.Vector{
		p.rotation.Mulv(p.meshes[index].GetVertex(0)).Add(p.position),
		p.rotation.Mulv(p.meshes[index].GetVertex(1)).Add(p.position),
		p.rotation.Mulv(p.meshes[index].GetVertex(2)).Add(p.position)})

	return mesh
}

// Meshes returns the number of meshes
func (p *Part) Meshes() int {
	return len(p.meshes)
}

// Box is a simple example that implements the Part interface
type Box struct {
	Part
	width  float32
	depth  float32
	height float32
}

// NewBox creates a simple Box part
func NewBox(width float32, depth float32, height float32) Box {
	var box Box

	// No negative stuff
	if width <= 0.0 || depth <= 0.0 || height <= 0.0 {
		log.Fatalf("Model.NewBox: must have positive sizes, got (w:%f, d:%f, h:%f)", width, depth, height)
	}

	// Default position at the origin
	box.position = vector.ZeroVector(3, reflect.Float32)

	// Default rotation at nothing
	box.rotation = matrix.UnitMatrix(3, 3, reflect.Float32)

	// Corner points
	fbl := vector.FilledVector([]float32{-width / 2.0, 0.0, -depth / 2.0})
	bbl := vector.FilledVector([]float32{-width / 2.0, 0.0, depth / 2.0})
	bbr := vector.FilledVector([]float32{width / 2.0, 0.0, depth / 2.0})
	fbr := vector.FilledVector([]float32{width / 2.0, 0.0, -depth / 2.0})
	ftl := vector.FilledVector([]float32{-width / 2.0, height, -depth / 2.0})
	btl := vector.FilledVector([]float32{-width / 2.0, height, depth / 2.0})
	btr := vector.FilledVector([]float32{width / 2.0, height, depth / 2.0})
	ftr := vector.FilledVector([]float32{width / 2.0, height, -depth / 2.0})

	// Meshes from the points
	box.meshes = []Mesh{
		// Bottom
		NewMesh([]vector.Vector{fbl, bbl, fbr}),
		NewMesh([]vector.Vector{bbl, bbr, fbr}),
		// Top
		NewMesh([]vector.Vector{ftl, btl, ftr}),
		NewMesh([]vector.Vector{btl, btr, ftr}),
		// Left
		NewMesh([]vector.Vector{fbl, bbl, btl}),
		NewMesh([]vector.Vector{fbl, ftl, btl}),
		// Right
		NewMesh([]vector.Vector{fbr, bbr, btr}),
		NewMesh([]vector.Vector{fbr, ftr, btr}),
		// Front
		NewMesh([]vector.Vector{fbl, ftl, ftr}),
		NewMesh([]vector.Vector{fbl, fbr, ftr}),
		// Back
		NewMesh([]vector.Vector{bbl, btl, btr}),
		NewMesh([]vector.Vector{bbl, bbr, btr}),
	}

	box.width = width
	box.depth = depth
	box.height = height

	return box
}
