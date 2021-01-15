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
	dimension int
	kind      reflect.Kind
	cells     interface{}
}

func genericZeroVector(dimension int, kind reflect.Kind) Vector {
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
		log.Panicf("genericVector.genericZeroVector: Unknown Kind for a Vector: %v\n", kind)
	}

	return genericVector{dimension, kind, cells}
}

// Internal 'rand' function
func genericRandomVector(dimension int, kind reflect.Kind) Vector {

	v := genericZeroVector(dimension, kind).(genericVector)
	for i := 0; i < v.dimension; i++ {
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
			log.Panicf("genericVector.genericRandomVector: Unknown Kind for a Vector: %v\n", kind)
		}
	}

	return v
}

func genericFilledVector(values interface{}) Vector {
	k := reflect.TypeOf(values).Kind()
	if k != reflect.Array && k != reflect.Slice {
		log.Fatalf("genericFilledVector: expected an array or slice, got %v", reflect.TypeOf(values).Kind())
	}

	source := reflect.ValueOf(values)
	if source.Len() == 0 {
		log.Fatalf("genericFilledVector: cannot create from an empty array")
	}

	// Fill with the content
	v := genericZeroVector(source.Len(), source.Index(0).Kind())
	for i := 0; i < v.Len(); i++ {
		v.Set(i, source.Index(i).Interface())
	}

	return v
}

// Unit provides a vector with length 1 in the direction of the vector
func (v genericVector) Unit() Vector {
	if v.Kind() != reflect.Float32 && v.Kind() != reflect.Float64 {
		log.Fatalf("genericVector.Unit: only supported for Float32 and FLoat64 vectors")
	}

	switch v.Kind() {
	case reflect.Float32:
		return v.Divs(float32(v.Abs()))
	case reflect.Float64:
		return v.Divs(v.Abs())
	}

	return genericZeroVector(v.Len(), v.Kind())
}

// Abs provides the euclidian lenth of a vector
// That's how mathematicians specify the absolute value of a vector
func (v genericVector) Abs() float64 {

	l := 0.0
	for i := 0; i < v.Len(); i++ {
		switch v.Kind() {
		case reflect.Int:
			l += math.Pow(float64(v.Get(i).(int)), 2)
		case reflect.Int8:
			l += math.Pow(float64(v.Get(i).(int8)), 2)
		case reflect.Int16:
			l += math.Pow(float64(v.Get(i).(int16)), 2)
		case reflect.Int32:
			l += math.Pow(float64(v.Get(i).(int32)), 2)
		case reflect.Int64:
			l += math.Pow(float64(v.Get(i).(int64)), 2)
		case reflect.Uint:
			l += math.Pow(float64(v.Get(i).(uint)), 2)
		case reflect.Uint8:
			l += math.Pow(float64(v.Get(i).(uint8)), 2)
		case reflect.Uint16:
			l += math.Pow(float64(v.Get(i).(uint16)), 2)
		case reflect.Uint32:
			l += math.Pow(float64(v.Get(i).(uint32)), 2)
		case reflect.Uint64:
			l += math.Pow(float64(v.Get(i).(uint64)), 2)
		case reflect.Float32:
			l += math.Pow(float64(v.Get(i).(float32)), 2)
		case reflect.Float64:
			l += math.Pow(v.Get(i).(float64), 2)
		default:
			log.Panicf("genericVector.Abs: Unknown Kind for a Vector: %v\n", v.Kind())
		}
	}

	return math.Sqrt(l)
}

// Cbd provides the city-block-distance length of a vector
// func (v genericVector) Cbd() interface{} {
// 	var l float64

// 	for i := 0; i < v.Len(); i++ {
// 		l += math.Abs(float64(v.Get(i)))
// 	}

// 	// TODO: the return value should be the same type as the vector?
// 	return l
// }

// Add substracts one vector from another
func (v genericVector) Add(w Vector) Vector {
	if w.Kind() != v.Kind() {
		log.Fatalf("genericVector.Add: kinds %v and %v do not match", v.Kind(), w.Kind())
	}
	if w.Len() != v.Len() {
		log.Fatalf("genericVector.Add: dimensions %d and %d do not match", v.Len(), w.Len())
	}

	r := genericZeroVector(v.Len(), v.Kind())
	for i := 0; i < r.Len(); i++ {
		switch v.Kind() {
		case reflect.Int:
			r.Set(i, v.Get(i).(int)+w.Get(i).(int))
		case reflect.Int8:
			r.Set(i, v.Get(i).(int8)+w.Get(i).(int8))
		case reflect.Int16:
			r.Set(i, v.Get(i).(int16)+w.Get(i).(int16))
		case reflect.Int32:
			r.Set(i, v.Get(i).(int32)+w.Get(i).(int32))
		case reflect.Int64:
			r.Set(i, v.Get(i).(int64)+w.Get(i).(int64))
		case reflect.Uint:
			r.Set(i, v.Get(i).(uint)+w.Get(i).(uint))
		case reflect.Uint8:
			r.Set(i, v.Get(i).(uint8)+w.Get(i).(uint8))
		case reflect.Uint16:
			r.Set(i, v.Get(i).(uint16)+w.Get(i).(uint16))
		case reflect.Uint32:
			r.Set(i, v.Get(i).(uint32)+w.Get(i).(uint32))
		case reflect.Uint64:
			r.Set(i, v.Get(i).(uint64)+w.Get(i).(uint64))
		case reflect.Float32:
			r.Set(i, v.Get(i).(float32)+w.Get(i).(float32))
		case reflect.Float64:
			r.Set(i, v.Get(i).(float64)+w.Get(i).(float64))
		default:
			log.Panicf("genericVector.Add: Unknown Kind for a Vector: %v\n", v.Kind())
		}
	}

	return r
}

// Sub substracts one vector from another
func (v genericVector) Sub(w Vector) Vector {
	if w.Kind() != v.Kind() {
		log.Fatalf("genericVector.Sub: kinds %v and %v do not match", v.Kind(), w.Kind())
	}
	if w.Len() != v.Len() {
		log.Fatalf("genericVector.Sub: dimensions %d and %d do not match", v.Len(), w.Len())
	}

	r := genericZeroVector(v.Len(), v.Kind())
	for i := 0; i < r.Len(); i++ {
		switch v.Kind() {
		case reflect.Int:
			r.Set(i, v.Get(i).(int)-w.Get(i).(int))
		case reflect.Int8:
			r.Set(i, v.Get(i).(int8)-w.Get(i).(int8))
		case reflect.Int16:
			r.Set(i, v.Get(i).(int16)-w.Get(i).(int16))
		case reflect.Int32:
			r.Set(i, v.Get(i).(int32)-w.Get(i).(int32))
		case reflect.Int64:
			r.Set(i, v.Get(i).(int64)-w.Get(i).(int64))
		case reflect.Uint:
			r.Set(i, v.Get(i).(uint)-w.Get(i).(uint))
		case reflect.Uint8:
			r.Set(i, v.Get(i).(uint8)-w.Get(i).(uint8))
		case reflect.Uint16:
			r.Set(i, v.Get(i).(uint16)-w.Get(i).(uint16))
		case reflect.Uint32:
			r.Set(i, v.Get(i).(uint32)-w.Get(i).(uint32))
		case reflect.Uint64:
			r.Set(i, v.Get(i).(uint64)-w.Get(i).(uint64))
		case reflect.Float32:
			r.Set(i, v.Get(i).(float32)-w.Get(i).(float32))
		case reflect.Float64:
			r.Set(i, v.Get(i).(float64)-w.Get(i).(float64))
		default:
			log.Panicf("genericVector.Sub: Unknown Kind for a Vector: %v\n", v.Kind())
		}
	}

	return r
}

// Min provides a vector containing the smallest elements of both vectors
// func (v genericVector) Min(w Vector) Vector {
// 	if v.Len() != w.Len() {
// 		log.Fatalf("Vector.Min: dimensions of vectors must be the same")
// 	}

// 	r := Zero(v.Len())
// 	for i := 0; i < r.Len(); i++ {
// 		r = r.Set(i, math.Min(v.Get(i).(float64), w.Get(i).(float64)))
// 	}

// 	return r
// }

// MinD provides the index of the smallest value in the vector
// func (v genericVector) MinD() int {

// 	r := 0
// 	for i := 0; i < v.Len(); i++ {
// 		if v.Get(i).(float64) < v.Get(r).(float64) {
// 			r = i
// 		}
// 	}

// 	return r
// }

// Max set every element of the resulting vector to the highest option
// It provides the first occurence if there are more dimensions with this value
// func (v genericVector) Max(w Vector) Vector {
// 	if v.Len() != w.Len() {
// 		log.Fatalf("Vector.Max: dimensions of vectors must be the same")
// 	}

// 	r := Zero(v.Len())
// 	for i := 0; i < r.Len(); i++ {
// 		r = r.Set(i, math.Max(v.Get(i).(float64), w.Get(i).(float64)))
// 	}

// 	return r
// }

// MaxD provides the index of the largest value in the vector
// It provides the first occurence if there are more dimensions with this value
// func (v genericVector) MaxD() int {

// 	r := 0
// 	for i := 0; i < v.Len(); i++ {
// 		if v.Get(i).(float64) > v.Get(r).(float64) {
// 			r = i
// 		}
// 	}

// 	return r
// }

// Muls multiplies a vector by a scalar.
// func (v genericVector) Muls(s interface{}) Vector {

// 	r := Zero(v.Len())
// 	for i := 0; i < r.Len(); i++ {
// 		r = r.Set(i, v.Get(i).(float64)*s.(float64))
// 	}

// 	return r
// }

// Divs divides a vector by a scalar.
func (v genericVector) Divs(s interface{}) Vector {

	if reflect.TypeOf(s).Kind() != v.Kind() {
		log.Fatalf("genericVector.Divs: Scalar Type %v doesn't match vector type %v", reflect.TypeOf(s).Kind(), v.Kind())
	}

	r := genericZeroVector(v.Len(), v.Kind())
	for i := 0; i < r.Len(); i++ {
		switch v.Kind() {
		case reflect.Int:
			r.Set(i, v.Get(i).(int)/s.(int))
		case reflect.Int8:
			r.Set(i, v.Get(i).(int8)/s.(int8))
		case reflect.Int16:
			r.Set(i, v.Get(i).(int16)/s.(int16))
		case reflect.Int32:
			r.Set(i, v.Get(i).(int32)/s.(int32))
		case reflect.Int64:
			r.Set(i, v.Get(i).(int64)/s.(int64))
		case reflect.Uint:
			r.Set(i, v.Get(i).(uint)/s.(uint))
		case reflect.Uint8:
			r.Set(i, v.Get(i).(uint8)/s.(uint8))
		case reflect.Uint16:
			r.Set(i, v.Get(i).(uint16)/s.(uint16))
		case reflect.Uint32:
			r.Set(i, v.Get(i).(uint32)/s.(uint32))
		case reflect.Uint64:
			r.Set(i, v.Get(i).(uint64)/s.(uint64))
		case reflect.Float32:
			r.Set(i, v.Get(i).(float32)/s.(float32))
		case reflect.Float64:
			r.Set(i, v.Get(i).(float64)/s.(float64))
		}
	}

	return r
}

// Kind retrieves the kind of values stored
func (v genericVector) Kind() reflect.Kind {
	return v.kind
}

// Len retrieves the number of dimensions
func (v genericVector) Len() int {
	return v.dimension
}

func (v genericVector) Equal(w Vector) bool {
	if w.Kind() != v.Kind() {
		log.Fatalf("genericVector.Equal: kinds %v and %v do not match", v.Kind(), w.Kind())
	}
	if w.Len() != v.Len() {
		log.Fatalf("genericVector.Equal: dimensions %d and %d do not match", v.Len(), w.Len())
	}

	equal := true
	for i := 0; (i < v.Len()) && equal; i++ {
		switch v.Kind() {
		case reflect.Int:
			equal = (v.Get(i).(int) == w.Get(i).(int))
		case reflect.Int8:
			equal = (v.Get(i).(int8) == w.Get(i).(int8))
		case reflect.Int16:
			equal = (v.Get(i).(int16) == w.Get(i).(int16))
		case reflect.Int32:
			equal = (v.Get(i).(int32) == w.Get(i).(int32))
		case reflect.Int64:
			equal = (v.Get(i).(int64) == w.Get(i).(int64))
		case reflect.Uint:
			equal = (v.Get(i).(uint) == w.Get(i).(uint))
		case reflect.Uint8:
			equal = (v.Get(i).(uint8) == w.Get(i).(uint8))
		case reflect.Uint16:
			equal = (v.Get(i).(uint16) == w.Get(i).(uint16))
		case reflect.Uint32:
			equal = (v.Get(i).(uint32) == w.Get(i).(uint32))
		case reflect.Uint64:
			equal = (v.Get(i).(uint64) == w.Get(i).(uint64))
		case reflect.Float32:
			equal = (v.Get(i).(float32) == w.Get(i).(float32))
		case reflect.Float64:
			equal = (v.Get(i).(float64) == w.Get(i).(float64))
		default:
			log.Panicf("genericVector.Divs: Unknown Kind for a Vector: %v\n", v.Kind())
		}
	}

	return equal
}

// Get retrieves the value of a single cell
func (v genericVector) Get(i int) interface{} {
	if i >= v.dimension || i >= math.MaxInt32 {
		log.Fatalf("genericVector.Set: Index %d out of bounds expected < %d", i, v.dimension)
	}
	return reflect.ValueOf(v.cells).Index(int(i)).Interface()
}

// Set changes the value of a single cell
// TODO: contemplate if 'Set' shouldn't be a private method
func (v genericVector) Set(i int, value interface{}) Vector {
	if i >= v.dimension || i >= math.MaxInt32 {
		log.Fatalf("genericVector.Set: Index %d out of bounds expected < %d", i, v.dimension)
	}
	if reflect.ValueOf(value).Type() != reflect.ValueOf(v.cells).Index(0).Type() {
		log.Fatalf("genericVector.Set: wrong value type %v expected %v", reflect.TypeOf(value), reflect.ValueOf(v.cells).Index(0).Type())
	}
	reflect.ValueOf(v.cells).Index(i).Set(reflect.ValueOf(value))
	return v
}

// String() implements the Stringer interface
func (v genericVector) String() string {
	var s strings.Builder

	s.WriteString("[")
	for i := 0; i < v.dimension; i++ {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(fmt.Sprintf("%v", v.Get(i)))
	}
	s.WriteString("]")

	return s.String()
}
