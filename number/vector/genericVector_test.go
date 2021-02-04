package vector

import (
	"reflect"
	"testing"
)

func Test_GenericZeroVector(t *testing.T) {
	// Simple 3D integer vector
	v := genericZeroVector(3, reflect.Int)
	for i := 0; i < v.Len(); i++ {
		if v.Get(i).(int) != 0 {
			t.Errorf("v[0] = %d, want 0", v.Get(i).(int))
		}
	}
}

func Test_NewVector(t *testing.T) {
	// Simple 3D integer vector
	v := NewVector([]int{1, 2, 3})
	if v.Get(0).(int) != 1 {
		t.Errorf("v[0] = %d, want 1", v.Get(0).(int))
	}
	if v.Get(1).(int) != 2 {
		t.Errorf("v[0] = %d, want 2", v.Get(0).(int))
	}
	if v.Get(2).(int) != 3 {
		t.Errorf("v[0] = %d, want 3", v.Get(0).(int))
	}
}

func Test_GenericUnit(t *testing.T) {
	// See if we get '1.0' for each of the axis
	v1 := NewVector([]float32{4.0, 0.0, 0.0})
	w1 := v1.Unit()
	if w1.Get(0).(float32) != 1.0 {
		t.Errorf("w[0] = %f, want 1.0", w1.Get(0).(float32))
	}
}

func Test_GenericSub(t *testing.T) {
	// See if we get the additive identity vector for ints
	v1 := NewVector([]int{1, 2, 3})
	w1 := v1.Sub(v1)
	if w1.Get(0).(int) != 0 {
		t.Errorf("w[0] = %d, want 0", w1.Get(0).(int))
	}
	if w1.Get(1).(int) != 0 {
		t.Errorf("w[1] = %d, want 0", w1.Get(1).(int))
	}
	if w1.Get(2).(int) != 0 {
		t.Errorf("w[2] = %d, want 0", w1.Get(2).(int))
	}

	// See if we get the additive identity vector for floats
	v2 := NewVector([]float32{1.0, 2.2, 3.3})
	w2 := v2.Sub(v2)
	if w2.Get(0).(float32) != 0.0 {
		t.Errorf("w[0] = %f, want 0.0", w2.Get(0).(float32))
	}
	if w2.Get(1).(float32) != 0.0 {
		t.Errorf("w[1] = %f, want 0.0", w2.Get(1).(float32))
	}
	if w2.Get(2).(float32) != 0.0 {
		t.Errorf("w[2] = %f, want 0.0", w2.Get(2).(float32))
	}

}

func Test_GenericDivs(t *testing.T) {
	// See if we can do an identity div for ints
	v1 := NewVector([]int{1, 2, 3})
	w1 := v1.Divs(int(1))
	if !w1.Equal(v1) {
		t.Errorf("%v / 1 --> %v, expected %v", v1, w1, v1)
	}

	// See if we get the identity divs for floats
	v2 := NewVector([]float32{1.1, 2.2, 3.3})
	w2 := v2.Divs(float32(1.0))
	if !w2.Equal(v2) {
		t.Errorf("%v / 1.0 --> %v, expected %v", v2, w2, v2)
	}
}

func Test_GenericEqual(t *testing.T) {
	// Test zero and identity equality for an int
	v1 := NewVector([]int{1, 2, 3})
	if v1.Equal(genericZeroVector(3, reflect.Int)) {
		t.Errorf("%v == %v", v1, genericZeroVector(3, reflect.Int))
	}
	if !v1.Equal(v1) {
		t.Errorf("%v != %v", v1, v1)
	}

	// Test zero identity equality for a float
	v2 := NewVector([]float32{1.1, 2.2, 3.3})
	if v2.Equal(genericZeroVector(3, reflect.Float32)) {
		t.Errorf("%v == %v", v2, genericZeroVector(3, reflect.Float32))
	}
	if !v2.Equal(v2) {
		t.Errorf("%v != %v", v2, v2)
	}

}
