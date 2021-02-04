package matrix

import (
	"fmt"
	"reflect"

	"../vector"
)

// Matrix interface allows to have specific types for various 'standard'
type Matrix interface {
	Mulv(v vector.Vector) vector.Vector
	Mulm(n Matrix) Matrix
	Kind() reflect.Kind
	Rows() int
	Cols() int
	Get(row int, col int) interface{}
	Set(row int, col int, value interface{})
	Equal(n Matrix) bool
	fmt.Stringer
}

// ZeroMatrix creates a matrix 'rows' high and 'cols' wide
func ZeroMatrix(rows int, cols int, kind reflect.Kind) Matrix {
	return genericZeroMatrix(rows, cols, kind)
}

// UnitMatrix creates a matrix with all ones across the main diagonal
func UnitMatrix(rows int, cols int, kind reflect.Kind) Matrix {
	return genericUnitMatrix(rows, cols, kind)
}

// NewMatrix creates a matrix based on a number of values
func NewMatrix(values interface{}) Matrix {
	return genericNewMatrix(values)
}
