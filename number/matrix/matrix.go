package matrix

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strings"
)

// Matrix interface allows to have specific types for various 'standard'
type Matrix interface {
	Rows() uint
	Cols() uint
	Get(row uint, col uint) interface{}
	Set(row uint, col uint, value interface{})
	fmt.Stringer
}

// The genericMatrix is an implementation of the Matrix interface which uses reflect to work with (almost) any type.
// We can specialize when runtime performance is painfull enough to warrant such.
type genericMatrix struct {
	rows   uint
	cols   uint
	kind   reflect.Kind
	values interface{}
}

func (m genericMatrix) Rows() uint {
	return m.rows
}

func (m genericMatrix) Cols() uint {
	return m.cols
}

func (m genericMatrix) Get(row uint, col uint) interface{} {
	if row >= m.rows || col >= m.cols {
		log.Panicf("genericMatrix.Get: index (%d, %d) out of bounds, expected(<%d, <%d)", row, col, m.rows, m.cols)
	}
	return reflect.ValueOf(m.values).Index(int(row*m.cols + col))
}

func (m genericMatrix) Set(row uint, col uint, value interface{}) {
	if row >= m.rows || col >= m.cols {
		log.Panicf("genericMatrix.Set: index (%d, %d) out of bounds, expected(<%d, <%d)", row, col, m.rows, m.cols)
	}
	if reflect.ValueOf(value).Type() != reflect.ValueOf(m.values).Index(0).Type() {
		log.Panicf("genericMatrix.Set: wrong value type %v, expected %v", reflect.ValueOf(value).Type(), reflect.ValueOf(m.values).Index(0).Type())
	}
	reflect.ValueOf(m.values).Index(int(row*m.cols + col)).Set(reflect.ValueOf(value))
}

func (m genericMatrix) String() string {
	var sb strings.Builder
	for r := uint(0); r < m.Rows(); r++ {
		for c := uint(0); c < m.Cols(); c++ {
			sb.WriteString(fmt.Sprintf("%v ", m.Get(r, c)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// NewMatrix creates a matrix 'rows' high and 'cols' wide
func NewMatrix(rows uint, cols uint, kind reflect.Kind) Matrix {
	// Check the maximum addressable space since we map a uint * uint onto an int
	if rows*cols > math.MaxInt64 {
		log.Panicf("NewMatrix: Size Overflow (%d, %d)", rows, cols)
	}

	// We need to create the array of the actual type in order to
	// have something addressable in GoLang
	var values interface{}

	switch kind {
	case reflect.Int:
		values = make([]int, rows*cols)
	case reflect.Int8:
		values = make([]int8, rows*cols)
	case reflect.Int16:
		values = make([]int16, rows*cols)
	case reflect.Int32:
		values = make([]int32, rows*cols)
	case reflect.Int64:
		values = make([]int64, rows*cols)
	case reflect.Uint:
		values = make([]uint, rows*cols)
	case reflect.Uint8:
		values = make([]uint8, rows*cols)
	case reflect.Uint16:
		values = make([]uint16, rows*cols)
	case reflect.Uint32:
		values = make([]uint32, rows*cols)
	case reflect.Uint64:
		values = make([]uint64, rows*cols)
	case reflect.Float32:
		values = make([]float32, rows*cols)
	case reflect.Float64:
		values = make([]float64, rows*cols)
	case reflect.Complex64:
		values = make([]complex64, rows*cols)
	case reflect.Complex128:
		values = make([]complex128, rows*cols)
	default:
		log.Panicf("Unknown Kind for a Matrix: %v\n", kind)
	}

	return genericMatrix{rows, cols, kind, values}
}
