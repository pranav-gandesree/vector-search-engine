package main

import (
	"fmt"
	"vector-search/index"
	"vector-search/index/flat"
	"vector-search/vector"
)

func main() {
	idx, _ := flat.NewIndex(2)

	record1 := index.Record{
		Id:     "A",
		Text:   "Database indexingg",
		Vector: []float64{1, 0},
	}

	record2 := index.Record{
		Id:     "B",
		Text:   "vector search",
		Vector: []float64{0, 1},
	}

	record3 := index.Record{
		Id:     "C",
		Text:   "vector search engine",
		Vector: []float64{-1, 0},
	}

	record4 := index.Record{
		Id:     "D",
		Text:   "distributed systems",
		Vector: []float64{1, 1},
	}

	idx.Add(record1)
	idx.Add(record2)
	idx.Add(record3)
	idx.Add(record4)

	query := vector.Vector{Values: []float64{0, 1}}

	results, err := idx.Search(query, 3)

	if err != nil {
		fmt.Println("Search error:", err)
		return
	}

	for i, result := range results {
		fmt.Printf("Search result %d: %+v\n", i+1, result)
	}
}
