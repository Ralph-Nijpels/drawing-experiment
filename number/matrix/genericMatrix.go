package matrix

import (
	"fmt"
	"log"
	"math"
	"math/rand"
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

func genericUnitMatrix(rows int, cols int, kind reflect.Kind) Matrix {
	matrix := genericZeroMatrix(rows, cols, kind)

	var value interface{}

	// Get the generic value for 1
	switch kind {
	case reflect.Int:
		value = int(1)
	case reflect.Int8:
		value = int8(1)
	case reflect.Int16:
		value = int16(1)
	case reflect.Int32:
		value = int32(1)
	case reflect.Int64:
		value = int64(1)
	case reflect.Uint:
		value = uint(1)
	case reflect.Uint8:
		value = uint8(1)
	case reflect.Uint16:
		value = uint16(1)
	case reflect.Uint32:
		value = uint32(1)
	case reflect.Uint64:
		value = uint64(1)
	case reflect.Float32:
		value = float32(1.0)
	case reflect.Float64:
		value = float64(1.0)
	}

	// run the diagonal
	for i := 0; i < rows && i < cols; i++ {
		matrix.Set(i, i, value)
	}

	return matrix
}

func genericRandomMatrix(rows int, cols int, kind reflect.Kind) Matrix {
	matrix := genericZeroMatrix(rows, cols, kind)

	for r := 0; r < matrix.Rows(); r++ {
		for c := 0; c < matrix.Cols(); c++ {
			switch kind {
			case reflect.Int:
				matrix.Set(r, c, rand.Int()) // TODO: we're not covering negatives
			case reflect.Int8:
				matrix.Set(r, c, int8(rand.Int())) // TODO: we're not covering negatives
			case reflect.Int16:
				matrix.Set(r, c, int16(rand.Int())) // TODO: we're not covering negatives
			case reflect.Int32:
				matrix.Set(r, c, int32(rand.Int31())) // TODO: we're not covering negatives
			case reflect.Int64:
				matrix.Set(r, c, int64(rand.Int63())) // TODO: we're not covering negatives
			case reflect.Uint:
				matrix.Set(r, c, uint(rand.Int())) // TODO: we're not covering negatives
			case reflect.Uint8:
				matrix.Set(r, c, uint8(rand.Int()))
			case reflect.Uint16:
				matrix.Set(r, c, uint16(rand.Int()))
			case reflect.Uint32:
				matrix.Set(r, c, uint32(rand.Int31())) // TODO: we're one bit short!
			case reflect.Uint64:
				matrix.Set(r, c, uint64(rand.Int63())) // TODO: we're one bit short!
			case reflect.Float32:
				matrix.Set(r, c, float32(rand.Float32()))
			case reflect.Float64:
				matrix.Set(r, c, float64(rand.Float64()))
			}
		}
	}

	return matrix
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

func (m genericMatrix) Mulm(n Matrix) Matrix {
	// Validity checks
	if m.Kind() != n.Kind() {
		log.Fatalf("genericMatrix.Mulm: expected matrix type %v, got %v", m.Kind(), n.Kind())
	}
	if m.Cols() != n.Rows() {
		log.Fatalf("genericMatrix.Mulm: expected matrix with %d rows, got %d", m.Cols(), n.Rows())
	}

	// Matrix multiplication bit
	result := ZeroMatrix(m.Rows(), n.Cols(), m.Kind())
	for cn := 0; cn < n.Cols(); cn++ { // walk the columns of the right-hand side like it is a list of vectors
		for rn := 0; rn < n.Rows(); rn++ { // take each row of these vectors
			for rm := 0; rm < m.Rows(); rm++ { // multiply & add with the corresponding rows
				switch result.Kind() {
				case reflect.Int:
					result.Set(rm, cn, result.Get(rm, cn).(int)+m.Get(rm, rn).(int)*n.Get(rn, cn).(int))
				case reflect.Int8:
					result.Set(rm, cn, result.Get(rm, cn).(int8)+m.Get(rm, rn).(int8)*n.Get(rn, cn).(int8))
				case reflect.Int16:
					result.Set(rm, cn, result.Get(rm, cn).(int16)+m.Get(rm, rn).(int16)*n.Get(rn, cn).(int16))
				case reflect.Int32:
					result.Set(rm, cn, result.Get(rm, cn).(int32)+m.Get(rm, rn).(int32)*n.Get(rn, cn).(int32))
				case reflect.Int64:
					result.Set(rm, cn, result.Get(rm, cn).(int64)+m.Get(rm, rn).(int64)*n.Get(rn, cn).(int64))
				case reflect.Uint:
					result.Set(rm, cn, result.Get(rm, cn).(uint)+m.Get(rm, rn).(uint)*n.Get(rn, cn).(uint))
				case reflect.Uint8:
					result.Set(rm, cn, result.Get(rm, cn).(uint8)+m.Get(rm, rn).(uint8)*n.Get(rn, cn).(uint8))
				case reflect.Uint16:
					result.Set(rm, cn, result.Get(rm, cn).(uint16)+m.Get(rm, rn).(uint16)*n.Get(rn, cn).(uint16))
				case reflect.Uint32:
					result.Set(rm, cn, result.Get(rm, cn).(uint32)+m.Get(rm, rn).(uint32)*n.Get(rn, cn).(uint32))
				case reflect.Uint64:
					result.Set(rm, cn, result.Get(rm, cn).(uint64)+m.Get(rm, rn).(uint64)*n.Get(rn, cn).(uint64))
				case reflect.Float32:
					result.Set(rm, cn, result.Get(rm, cn).(float32)+m.Get(rm, rn).(float32)*n.Get(rn, cn).(float32))
				case reflect.Float64:
					result.Set(rm, cn, result.Get(rm, cn).(float64)+m.Get(rm, rn).(float64)*n.Get(rn, cn).(float64))
				}
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

func (m genericMatrix) Equal(n Matrix) bool {
	if n.Kind() != m.Kind() {
		log.Fatalf("genericMatrix.Equal: kinds %v and %v do not match", n.Kind(), m.Kind())
	}
	if n.Rows() != m.Rows() || n.Cols() != m.Cols() {
		log.Fatalf("genericMatrix.Equal: dimensions (%d, %d) and (%d, %d) do not match", m.Rows(), m.Cols(), n.Rows(), n.Cols())
	}

	equal := true
	for r := 0; (r < m.Rows()) && equal; r++ {
		for c := 0; (c < m.Cols()) && equal; c++ {
			switch m.Kind() {
			case reflect.Int:
				equal = (m.Get(r,c).(int) == n.Get(r,c).(int))
			case reflect.Int8:
				equal = (m.Get(r,c).(int8) == n.Get(r,c).(int8))
			case reflect.Int16:
				equal = (m.Get(r,c).(int16) == n.Get(r,c).(int16))
			case reflect.Int32:
				equal = (m.Get(r,c).(int32) == n.Get(r,c).(int32))
			case reflect.Int64:
				equal = (m.Get(r,c).(int64) == n.Get(r,c).(int64))
			case reflect.Uint:
				equal = (m.Get(r,c).(uint) == n.Get(r,c).(uint))
			case reflect.Uint8:
				equal = (m.Get(r,c).(uint8) == n.Get(r,c).(uint8))
			case reflect.Uint16:
				equal = (m.Get(r,c).(uint16) == n.Get(r,c).(uint16))
			case reflect.Uint32:
				equal = (m.Get(r,c).(uint32) == n.Get(r,c).(uint32))
			case reflect.Uint64:
				equal = (m.Get(r,c).(uint64) == n.Get(r,c).(uint64))
			case reflect.Float32:
				equal = (m.Get(r,c).(float32) == n.Get(r,c).(float32))
			case reflect.Float64:
				equal = (m.Get(r,c).(float64) == n.Get(r,c).(float64))
			}
		}
	}

	return equal
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
