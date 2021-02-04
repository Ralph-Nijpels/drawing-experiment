package render

import (
	"fmt"
	"testing"

	"../number/vector"
)

func Test_Project(t *testing.T) {
	pos := vector.NewVector([]float32{1.0, 1.0, 100.0})
	dir := vector.NewVector([]float32{0.0, 0.0, 0.0})
	camera := NewCamera(pos, dir)

	t0 := vector.NewVector([]float32{0.0, 0.0, 0.0})
	r0 := camera.Project(t0)
	fmt.Printf("Camera Projection: %v --> %v\n", t0, r0)
}
