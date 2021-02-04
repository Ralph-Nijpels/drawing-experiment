package vector

import (
	"fmt"
	"reflect"
)

// Vector implements a simple mathematical vector
type Vector interface {
	Unit() Vector
	Abs() float64
	//	Cbd() interface{}
	Add(w Vector) Vector
	Sub(w Vector) Vector
	//	Min(w Vector) Vector
	//	MinD() int
	//	Max(w Vector) Vector
	//	MaxD() int
	Muls(s interface{}) Vector
	Divs(s interface{}) Vector
	Mulv(w Vector) float64
	Kind() reflect.Kind
	Len() int
	Get(i int) interface{}
	Set(i int, f interface{}) Vector
	Equal(w Vector) bool
	fmt.Stringer
}

// ZeroVector creates a vector of the requested size set to the origin
func ZeroVector(dimension int, kind reflect.Kind) Vector {
	//	if dimension == 3 {
	//		return zero3D()
	//	}
	return genericZeroVector(dimension, kind)
}

// RandomVector creates a vector of the requested size set to a random location
func RandomVector(dimension int, kind reflect.Kind) Vector {
	//	if dimension == 3 {
	//		return rand3D()
	//	}
	return genericRandomVector(dimension, kind)
}

// NewVector creates a vector based on a list of values
func NewVector(f interface{}) Vector {
	//	if len(f) == 3 {
	//		return make3D(f)
	//	}
	return genericNewVector(f)
}
