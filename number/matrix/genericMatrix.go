package matrix

import (
	"fmt"
	"log"
	"math"
	"reflect"
	"strings"

	"../vector"
)

// The genericMatrix is an implementation of the Matrix interface which uses reflect to work with (almost) any type.
// We can specialize when runtime performance is painfull enough to warrant such.
type genericMatrix struct {
	rows   int
	cols   int
	kind   reflect.Kind
	values interface{}
}

func genericZeroMatrix(rows int, cols int, kind reflect.Kind) Matrix {
	// Check the maximum addressable space since we map a uint * uint onto an int
	if rows*cols > math.MaxInt32 {
		log.Panicf("genericZeroMatrix: Size Overflow (%d, %d)", rows, cols)
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
	default:
		log.Panicf("Unknown Kind for a Matrix: %v\n", kind)
	}

	return genericMatrix{rows, cols, kind, values}
}

func genericFilledMatrix(values interface{}) Matrix {
	source := reflect.ValueOf(values)
	// Check row-level
	kind := source.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		log.Fatalf("genericFilledMatrix: expected an array or slice, got %v", kind)
	}
	rows := source.Len()
	if rows == 0 {
		log.Fatalf("genericFilledMatrix: cannot create with zero rows")
	}
	// Check col-level
	kind = source.Index(0).Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		log.Fatalf("genericFilledMatrix: expected an array or slice, got %v", kind)
	}
	cols := source.Index(0).Len()
	if cols == 0 {
		log.Fatalf("genericFilledMatrix: cannot create with zero cols")
	}

	// Fill with the content
	kind = source.Index(0).Index(0).Kind()
	matrix := genericZeroMatrix(rows, cols, kind)
	for r := 0; r < rows; r++ {
		row := source.Index(r)
		if row.Len() != cols {
			log.Fatalf("genericFilledMatrix: irregular row %d expected %d cols, got %d", r, cols, row.Len())
		}
		for c := 0; c < cols; c++ {
			matrix.Set(r, c, source.Index(r).Index(c).Interface())
		}
	}

	return matrix
}

func (m genericMatrix) Mulv(v vector.Vector) vector.Vector {
	// Validity checks
	if m.Kind() != v.Kind() {
		log.Fatalf("genericMatrix.Mulv: expected vector type %v, got %v", m.Kind(), v.Kind())
	}
	if m.Cols() != v.Len() {
		log.Fatalf("genericMatrix.Mulv: expected vector length %d, got %d", m.Cols(), v.Len())
	}

	result := vector.ZeroVector(m.rows, m.Kind())
	for r := 0; r < m.Rows(); r++ {
		for c := 0; c < m.Cols(); c++ {
			switch result.Kind() {
			case reflect.Int:
				result.Set(r, result.Get(r).(int)+m.Get(r, c).(int)*v.Get(c).(int))
			case reflect.Int8:
				result.Set(r, result.Get(r).(int8)+m.Get(r, c).(int8)*v.Get(c).(int8))
			case reflect.Int16:
				result.Set(r, result.Get(r).(int16)+m.Get(r, c).(int16)*v.Get(c).(int16))
			case reflect.Int32:
				result.Set(r, result.Get(r).(int32)+m.Get(r, c).(int32)*v.Get(c).(int32))
			case reflect.Int64:
				result.Set(r, result.Get(r).(int64)+m.Get(r, c).(int64)*v.Get(c).(int64))
			case reflect.Uint:
				result.Set(r, result.Get(r).(uint)+m.Get(r, c).(uint)*v.Get(c).(uint))
			case reflect.Uint8:
				result.Set(r, result.Get(r).(uint8)+m.Get(r, c).(uint8)*v.Get(c).(uint8))
			case reflect.Uint16:
				result.Set(r, result.Get(r).(uint16)+m.Get(r, c).(uint16)*v.Get(c).(uint16))
			case reflect.Uint32:
				result.Set(r, result.Get(r).(uint32)+m.Get(r, c).(uint32)*v.Get(c).(uint32))
			case reflect.Uint64:
				result.Set(r, result.Get(r).(uint64)+m.Get(r, c).(uint64)*v.Get(c).(uint64))
			case reflect.Float32:
				result.Set(r, result.Get(r).(float32)+m.Get(r, c).(float32)*v.Get(c).(float32))
			case reflect.Float64:
				result.Set(r, result.Get(r).(float64)+m.Get(r, c).(float64)*v.Get(c).(float64))
			}
		}
	}

	return result
}

func (m genericMatrix) Kind() reflect.Kind {
	return m.kind
}

func (m genericMatrix) Rows() int {
	return m.rows
}

func (m genericMatrix) Cols() int {
	return m.cols
}

func (m genericMatrix) Get(row int, col int) interface{} {
	if row >= m.rows || col >= m.cols {
		log.Panicf("genericMatrix.Get: index (%d, %d) out of bounds, expected(<%d, <%d)", row, col, m.rows, m.cols)
	}
	return reflect.ValueOf(m.values).Index(int(row*m.cols + col)).Interface()
}

func (m genericMatrix) Set(row int, col int, value interface{}) {
	if row >= m.rows || col >= m.cols {
		log.Panicf("genericMatrix.Set: index (%d, %d) out of bounds, expected(<%d, <%d)", row, col, m.rows, m.cols)
	}
	if reflect.ValueOf(value).Kind() != m.Kind() {
		log.Panicf("genericMatrix.Set: wrong value type %v, expected %v", reflect.ValueOf(value).Kind(), m.Kind())
	}
	reflect.ValueOf(m.values).Index(row*m.cols + col).Set(reflect.ValueOf(value))
}

func (m genericMatrix) String() string {
	var sb strings.Builder
	for r := 0; r < m.Rows(); r++ {
		for c := 0; c < m.Cols(); c++ {
			sb.WriteString(fmt.Sprintf("%v ", m.Get(r, c)))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
