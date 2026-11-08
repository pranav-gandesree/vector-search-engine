package main

import (
	"fmt"
	"vector-search/vector"
)

func main() {
	v1 := vector.Vector{Values: []float64{3e200, 4e200}}
	v2 := vector.Vector{Values: []float64{3e-200, 4e-200}}

	// result, err := v1.Add(v2)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println(result.Values)

	result1, _ := vector.Normalization(v1)

	fmt.Println(result1.Values)

	result2, _ := vector.DotProduct(v1, v2)

	fmt.Println(result2)

}
