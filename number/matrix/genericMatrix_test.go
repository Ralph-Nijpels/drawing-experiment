package matrix

import (
	"reflect"
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

func Test_GenericMulm(t *testing.T) {
	// Check the mul identity matrix for integers
	m0 := genericUnitMatrix(3, 3, reflect.Int)
	n0 := genericRandomMatrix(3, 4, reflect.Int)

	r0 := m0.Mulm(n0)
	if !r0.Equal(n0) {
		t.Errorf("Expected %v, got %v", n0, r0)
	}

	// The other way arround should be the same
	m1 := genericRandomMatrix(3, 3, reflect.Int)
	n1 := genericUnitMatrix(3, 3, reflect.Int)

	r1 := m1.Mulm(n1)
	if !r1.Equal(m1) {
		t.Errorf("Expected %v, got %v", m1, r1)
	}

	// Check the mul identity matrix for floats
	m2 := genericUnitMatrix(3, 3, reflect.Float32)
	n2 := genericRandomMatrix(3, 4, reflect.Float32)

	r2 := m2.Mulm(n2)
	if !r2.Equal(n2) {
		t.Errorf("Expected %v, got %v", n2, r2)
	}

	// The other way arround should be the same
	m3 := genericRandomMatrix(3, 3, reflect.Float32)
	n3 := genericUnitMatrix(3, 3, reflect.Float32)

	r3 := m3.Mulm(n3)
	if !r3.Equal(m3) {
		t.Errorf("Expected %v, got %v", m3, r3)
	}

}
