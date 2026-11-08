package vector

import (
	"fmt"
)

type Vector struct {
	Values []float64
}

func (v Vector) Add(other Vector) (Vector, error) {
	if len(v.Values) != len(other.Values) {
		return Vector{}, fmt.Errorf("vector dimensions must match")
	}

	result := make([]float64, len(v.Values))

	for i := range v.Values {
		result[i] = v.Values[i] + other.Values[i]
	}

	return Vector{Values: result}, nil
}
