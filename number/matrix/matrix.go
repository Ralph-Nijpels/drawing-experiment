package matrix

import (
	"fmt"
	"reflect"

	"../vector"
)

// Matrix interface allows to have specific types for various 'standard'
type Matrix interface {
	Mulv(v vector.Vector) vector.Vector
	Kind() reflect.Kind
	Rows() int
	Cols() int
	Get(row int, col int) interface{}
	Set(row int, col int, value interface{})
	fmt.Stringer
}

// ZeroMatrix creates a matrix 'rows' high and 'cols' wide
func ZeroMatrix(rows int, cols int, kind reflect.Kind) Matrix {
	return genericZeroMatrix(rows, cols, kind)
}

// FilledMatrix creates a matrix based on a number of values
func FilledMatrix(values interface{}) Matrix {
	return genericFilledMatrix(values)
}
