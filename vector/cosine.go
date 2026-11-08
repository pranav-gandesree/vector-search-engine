package vector

import (
	"fmt"
	"math"
)

// Implement normalization and cosine similarity here.
func Normalization(vector Vector) (Vector, error) {

	//check if vector is empty
	if len(vector.Values) == 0 {
		return Vector{}, fmt.Errorf("vector must not be empty")
	}

	// validate each component
	for i, value := range vector.Values {
		if math.IsNaN(value) {
			return Vector{}, fmt.Errorf("component %d is NaN", i)
		}

		if math.IsInf(value, 0) {
			return Vector{}, fmt.Errorf("component %d is infinite", i)
		}
	}

	//fing largest absolute value because squaring too large or too small values is risky
	var maxAbs float64

	for _, value := range vector.Values {
		absValue := math.Abs(value)

		if absValue > maxAbs {
			maxAbs = absValue
		}
	}
	//reject zero vector
	if maxAbs == 0 {
		return Vector{}, fmt.Errorf("cannot normalize zero vector")
	}

	var sum float64

	for _, value := range vector.Values {
		scaled := value / maxAbs
		sum += scaled * scaled
	}

	magnitude := math.Sqrt(sum)

	if magnitude == 0 {
		return Vector{}, fmt.Errorf("cannot normalize zero vector")
	}

	normalisedValues := make([]float64, len(vector.Values))

	for i, value := range vector.Values {
		normalisedValues[i] = (value / maxAbs) / magnitude
	}

	return Vector{Values: normalisedValues}, nil

}

func DotProduct(vector1, vector2 Vector) (float64, error) {
	if len(vector1.Values) != len(vector2.Values) {
		return 0, fmt.Errorf("vector dimensions must match")
	}

	if len(vector1.Values) == 0 || len(vector2.Values) == 0 {
		return 0, fmt.Errorf("vectors must not be empty")
	}

	var dotprod float64

	for i, value := range vector1.Values {

		if math.IsNaN(value) || math.IsInf(value, 0) {
			return 0, fmt.Errorf("vector1 component %d is not finite", i)
		}

		other := vector2.Values[i]

		if math.IsNaN(other) || math.IsInf(other, 0) {
			return 0, fmt.Errorf("vector2 component %d is not finite", i)
		}

		dotprod += value * other
	}

	if math.IsNaN(dotprod) || math.IsInf(dotprod, 0) {
		return 0, fmt.Errorf("dot product is not finite")
	}

	return dotprod, nil
}
