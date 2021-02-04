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
	shearing matrix.Matrix
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
		log.Fatalf("Part.SetRotation: expects 3D-Float32 vector, got %dD-%v", rotation.Len(), rotation.Kind())
	}
	// Make float64 values in Radians for math
	xAngle := float64(rotation.Get(0).(float32)) * math.Pi / 180.0
	yAngle := float64(rotation.Get(1).(float32)) * math.Pi / 180.0
	zAngle := float64(rotation.Get(2).(float32)) * math.Pi / 180.0
	// Create a rotation matrix for each axis
	xRotation := matrix.NewMatrix([][]float32{
		{1.0, 0.0, 0.0},
		{0.0, float32(math.Cos(xAngle)), float32(-math.Sin(xAngle))},
		{0.0, float32(math.Sin(xAngle)), float32(math.Cos(xAngle))},
	})
	yRotation := matrix.NewMatrix([][]float32{
		{float32(math.Cos(yAngle)), 0.0, float32(math.Sin(yAngle))},
		{0.0, 1.0, 0.0},
		{float32(-math.Sin(yAngle)), 0.0, float32(math.Cos(yAngle))},
	})
	zRotation := matrix.NewMatrix([][]float32{
		{float32(math.Cos(zAngle)), float32(-math.Sin(zAngle)), 0.0},
		{float32(math.Sin(zAngle)), float32(math.Cos(zAngle)), 0.0},
		{0.0, 0.0, 1.0},
	})
	// Combine the rotations
	p.rotation = zRotation.Mulm(yRotation.Mulm(xRotation))
}

// SetScale scales the part within it's own coordinate system
func (p *Part) SetScale(scale vector.Vector) {
	if scale.Len() != 3 && scale.Kind() != reflect.Float32 {
		log.Fatalf("Part.SetScale: expects 3D-Float32 vector, got %dD-%v", scale.Len(), scale.Kind())
	}
	// Create scaling matrix
	p.scaling = matrix.NewMatrix([][]float32{
		{scale.Get(0).(float32), 0.0, 0.0},
		{0.0, scale.Get(1).(float32), 0.0},
		{0.0, 0.0, scale.Get(2).(float32)},
	})
}

// SetShear shears the part within it's own coordinate system
func (p *Part) SetShear(sheer vector.Vector) {
	log.Fatalf("Part.SetShear: not yet implemented")
}

// GetMeshes returns a list of meshes for the entire part: scaled, sheared, rotated and positioned
func (p *Part) GetMeshes() []Mesh {

	// Set-up the translation matrix
	if p.position == nil {
		p.position = vector.ZeroVector(3, reflect.Float32)
	}
	if p.rotation == nil {
		p.rotation = matrix.UnitMatrix(3, 3, reflect.Float32)
	}
	if p.scaling == nil {
		p.scaling = matrix.UnitMatrix(3, 3, reflect.Float32)
	}
	if p.shearing == nil {
		p.shearing = matrix.UnitMatrix(3, 3, reflect.Float32)
	}
	translation := p.rotation.Mulm(p.shearing.Mulm(p.scaling))

	// TODO: process subparts

	// scale, shear, rotate and reposition
	result := make([]Mesh, len(p.meshes))
	for index, mesh := range p.meshes {
		this := NewMesh([]vector.Vector{
			translation.Mulv(mesh.GetVertex(0)).Add(p.position),
			translation.Mulv(mesh.GetVertex(1)).Add(p.position),
			translation.Mulv(mesh.GetVertex(2)).Add(p.position),
		})
		result[index] = this
	}

	return result
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

	// Corner points
	fbl := vector.NewVector([]float32{-width / 2.0, 0.0, -depth / 2.0})
	bbl := vector.NewVector([]float32{-width / 2.0, 0.0, depth / 2.0})
	bbr := vector.NewVector([]float32{width / 2.0, 0.0, depth / 2.0})
	fbr := vector.NewVector([]float32{width / 2.0, 0.0, -depth / 2.0})
	ftl := vector.NewVector([]float32{-width / 2.0, height, -depth / 2.0})
	btl := vector.NewVector([]float32{-width / 2.0, height, depth / 2.0})
	btr := vector.NewVector([]float32{width / 2.0, height, depth / 2.0})
	ftr := vector.NewVector([]float32{width / 2.0, height, -depth / 2.0})

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
