package vector

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"reflect"
	"strings"
)

// VectorDyn implements the dynamic (hence slow) version of the vector interface
type genericVector struct {
	dimension uint
	cells     interface{}
}

func genericZeroVector(dimension uint, kind reflect.Kind) Vector {
	var cells interface{}

	// Create the vector using the first element of the array as a template
	switch kind {
	case reflect.Int:
		cells = make([]int, dimension)
	case reflect.Int8:
		cells = make([]int8, dimension)
	case reflect.Int16:
		cells = make([]int16, dimension)
	case reflect.Int32:
		cells = make([]int32, dimension)
	case reflect.Int64:
		cells = make([]int64, dimension)
	case reflect.Uint:
		cells = make([]uint, dimension)
	case reflect.Uint8:
		cells = make([]uint8, dimension)
	case reflect.Uint16:
		cells = make([]uint16, dimension)
	case reflect.Uint32:
		cells = make([]uint32, dimension)
	case reflect.Uint64:
		cells = make([]uint64, dimension)
	case reflect.Float32:
		cells = make([]float32, dimension)
	case reflect.Float64:
		cells = make([]float64, dimension)
	default:
		log.Panicf("Unknown Kind for a Vector: %v\n", kind)
	}

	return genericVector{dimension, cells}
}

// Internal 'rand' function
func genericRandomVector(dimension uint, kind reflect.Kind) Vector {
	v := genericZeroVector(dimension, kind).(genericVector)
	for i := uint(0); i < v.dimension; i++ {
		switch kind {
		case reflect.Int:
			v.Set(i, rand.Int()) // TODO: we're not covering negatives
		case reflect.Int8:
			v.Set(i, int8(rand.Int())) // TODO: we're not covering negatives
		case reflect.Int16:
			v.Set(i, int16(rand.Int())) // TODO: we're not covering negatives
		case reflect.Int32:
			v.Set(i, int32(rand.Int31())) // TODO: we're not covering negatives
		case reflect.Int64:
			v.Set(i, int64(rand.Int63())) // TODO: we're not covering negatives
		case reflect.Uint:
			v.Set(i, uint(rand.Int())) // TODO: we're not covering negatives
		case reflect.Uint8:
			v.Set(i, uint8(rand.Int()))
		case reflect.Uint16:
			v.Set(i, uint16(rand.Int()))
		case reflect.Uint32:
			v.Set(i, uint32(rand.Int31())) // TODO: we're one bit short!
		case reflect.Uint64:
			v.Set(i, uint64(rand.Int63())) // TODO: we're one bit short!
		case reflect.Float32:
			v.Set(i, float32(rand.Float32()))
		case reflect.Float64:
			v.Set(i, float64(rand.Float64()))
		default:
			log.Panicf("Unknown Kind for a Vector: %v\n", kind)
		}
	}

	return v
}

func genericFilledVector(values []interface{}) Vector {

	if len(values) == 0 {
		log.Panicf("vector::genericVector cannot create from an empty array")
	}

	// Fill with the content
	v := genericZeroVector(uint(len(values)), reflect.TypeOf(values[0]).Kind())
	for i := uint(0); i < v.Len(); i++ {
		v.Set(i, values[i])
	}

	return v
}

/*
// Unit provides a vector with length 1 in the direction of the vector
func (v genericVector) Unit() Vector {
	return v.Divs(v.Abs())
}

// Abs provides the euclidian lenth of a vector
// That's how mathematicians specify the absolute value of a vector
func (v genericVector) Abs() interface{} {

	l := 0.0
	for i := 0; i < v.Len(); i++ {
		l += math.Pow(float64(v.Get(i)), 2)
	}

	return math.Sqrt(l)
}

// Cbd provides the city-block-distance length of a vector
func (v genericVector) Cbd() interface{} {
	var l float64

	for i := 0; i < v.Len(); i++ {
		l += math.Abs(float64(v.Get(i)))
	}

	// TODO: the return value should be the same type as the vector?
	return l
}

// Add substracts one vector from another
func (v genericVector) Add(w Vector) Vector {
	if v.Len() != w.Len() {
		log.Fatalf("Vector.Add: dimensions of vectors must be the same")
	}

	r := Zero(v.Len())
	for i := 0; i < r.Len(); i++ {
		r = r.Set(i, v.Get(i)+w.Get(i))
	}

	return r
}

// Sub substracts one vector from another
func (v genericVector) Sub(w Vector) Vector {
	if v.Len() != w.Len() {
		log.Fatalf("Vector.Sub: dimensions of vectors must be the same")
	}

	r := Zero(v.Len())
	for i := 0; i < r.Len(); i++ {
		r = r.Set(i, v.Get(i).(float64)-w.Get(i).(float64))
	}

	return r
}

// Min provides a vector containing the smallest elements of both vectors
func (v genericVector) Min(w Vector) Vector {
	if v.Len() != w.Len() {
		log.Fatalf("Vector.Min: dimensions of vectors must be the same")
	}

	r := Zero(v.Len())
	for i := 0; i < r.Len(); i++ {
		r = r.Set(i, math.Min(v.Get(i).(float64), w.Get(i).(float64)))
	}

	return r
}

// MinD provides the index of the smallest value in the vector
func (v genericVector) MinD() int {

	r := 0
	for i := 0; i < v.Len(); i++ {
		if v.Get(i).(float64) < v.Get(r).(float64) {
			r = i
		}
	}

	return r
}

// Max set every element of the resulting vector to the highest option
// It provides the first occurence if there are more dimensions with this value
func (v genericVector) Max(w Vector) Vector {
	if v.Len() != w.Len() {
		log.Fatalf("Vector.Max: dimensions of vectors must be the same")
	}

	r := Zero(v.Len())
	for i := 0; i < r.Len(); i++ {
		r = r.Set(i, math.Max(v.Get(i).(float64), w.Get(i).(float64)))
	}

	return r
}

// MaxD provides the index of the largest value in the vector
// It provides the first occurence if there are more dimensions with this value
func (v genericVector) MaxD() int {

	r := 0
	for i := 0; i < v.Len(); i++ {
		if v.Get(i).(float64) > v.Get(r).(float64) {
			r = i
		}
	}

	return r
}

// Muls multiplies a vector by a scalar.
func (v genericVector) Muls(s interface{}) Vector {

	r := Zero(v.Len())
	for i := 0; i < r.Len(); i++ {
		r = r.Set(i, v.Get(i).(float64)*s.(float64))
	}

	return r
}

// Divs divides a vector by a scalar.
func (v genericVector) Divs(s interface{}) Vector {

	r := Zero(v.Len())
	for i := 0; i < r.Len(); i++ {
		r = r.Set(i, v.Get(i).(float64)/s.(float64))
	}

	return r
}

*/
// Len retrieves the number of dimensions
func (v genericVector) Len() uint {
	return v.dimension
}

// Get retrieves the value of a single cell
func (v genericVector) Get(i uint) interface{} {
	if i >= v.dimension || i >= math.MaxInt32 {
		log.Fatalf("genericVector.Set: Index %d out of bounds expected < %d", i, v.dimension)
	}

	return reflect.ValueOf(v.cells).Index(int(i))
}

// Set changes the value of a single cell
// TODO: contemplate if 'Set' shouldn't be a private method
func (v genericVector) Set(i uint, value interface{}) Vector {
	if i >= v.dimension || i >= math.MaxInt32 {
		log.Fatalf("genericVector.Set: Index %d out of bounds expected < %d", i, v.dimension)
	}
	if reflect.ValueOf(value).Type() != reflect.ValueOf(v.cells).Index(0).Type() {
		log.Fatalf("genericVector.Set: wrong value type %v expected %v", reflect.TypeOf(value), reflect.ValueOf(v.cells).Index(0).Type())
	}
	reflect.ValueOf(v.cells).Index(int(i)).Set(reflect.ValueOf(value))
	return v
}

// String() implements the Stringer interface
func (v genericVector) String() string {
	var s strings.Builder

	s.WriteString("[")
	for i := uint(0); i < v.dimension; i++ {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprintf("%v", v.Get(i)))
	}
	s.WriteString("]")

	return s.String()
}
