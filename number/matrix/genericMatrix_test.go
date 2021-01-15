package matrix

import (
	"testing"

	"../vector"
)

func Test_GenericFilledMatrix(t *testing.T) {
	m := genericFilledMatrix([][]int{
		{1, 2},
		{3, 4},
	})
	if m.Get(0, 0).(int) != 1 {
		t.Errorf("Filled Matrix = %v", m)
	}
}

func Test_GenericMulv(t *testing.T) {
	// Check the mul identity matrix for integers
	m0 := genericFilledMatrix([][]int{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	})
	v := vector.FilledVector([]int{1, 2, 3})
	r := m0.Mulv(v)

	if !r.Equal(v) {
		t.Errorf("Expected %v, got %v", v, r)
	}

	// Check the mul selection matrix for integers
	m1 := genericFilledMatrix([][]int{
		{1, 0, 0},
		{1, 0, 0},
		{1, 0, 0},
	})
	r = m1.Mulv(v)
	if !r.Equal(vector.FilledVector([]int{1, 1, 1})) {
		t.Errorf("Expected %v, got %v", vector.FilledVector([]int{1, 1, 1}), r)
	}

	// Check the mul as a simple projection for integers
	m2 := genericFilledMatrix([][]int{
		{1, 0, 1},
		{0, 1, 1},
	})
	r = m2.Mulv(v)
	if !r.Equal(vector.FilledVector([]int{4, 5})) {
		t.Errorf("Expected %v, got %v", vector.FilledVector([]int{4, 5}), r)
	}
}
