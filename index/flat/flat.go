// Package flat will implement exact nearest-neighbor search.
package flat

import (
	"errors"
	"fmt"
	"sort"
	"vector-search/index"
	"vector-search/vector"
)

type Index struct {
	Dimension int
	Records   []index.Record
	IDs       map[string]struct{}
}

func NewIndex(dimension int) (*Index, error) {
	if dimension < 1 {
		return nil, errors.New("dimension must be at least 1")
	}

	return &Index{
		Dimension: dimension,
		Records:   make([]index.Record, 0),
		IDs:       make(map[string]struct{}),
	}, nil
}

func (i *Index) Add(record index.Record) error {
	if record.Id == "" {
		return fmt.Errorf("record Id must not be empty")
	}

	if _, exists := i.IDs[record.Id]; exists {
		return fmt.Errorf("This Id already exists %s", record.Id)
	}

	if len(record.Vector) != i.Dimension {
		return fmt.Errorf("new vector dimension must be %d", i.Dimension)
	}

	//normalising the vector
	normalisedVector, err := vector.Normalization(
		vector.Vector{Values: record.Vector},
	)

	if err != nil {
		return err
	}

	// store the normalized vector
	record.Vector = normalisedVector.Values

	// add the record
	i.Records = append(i.Records, record)

	// add ID to the set
	i.IDs[record.Id] = struct{}{}

	return nil
}

func (i *Index) Search(query vector.Vector, k int) ([]index.SearchResult, error) {
	if k <= 0 {
		return nil, fmt.Errorf("k must be greater than zero")
	}

	if len(query.Values) != i.Dimension {
		return nil, fmt.Errorf(
			"expected dimension %d, got %d", i.Dimension, len(query.Values),
		)
	}

	normalizedQuery, err := vector.Normalization(query)

	if err != nil {
		return nil, err
	}

	var scoredRecords []index.SearchResult

	for _, record := range i.Records {
		//stored vector is already normalised
		storedVector := vector.Vector{
			Values: record.Vector,
		}

		//now calculate the cosine similarity using dot product
		score, err := vector.DotProduct(normalizedQuery, storedVector)
		if err != nil {
			return nil, err
		}

		scoredRecords = append(scoredRecords, index.SearchResult{
			ItemId:          record.Id,
			SimilarityScore: score,
		})
	}

	// sort highest similarity first
	sort.Slice(scoredRecords, func(a, b int) bool {
		if scoredRecords[a].SimilarityScore == scoredRecords[b].SimilarityScore {
			return scoredRecords[a].ItemId < scoredRecords[b].ItemId
		}
		return scoredRecords[a].SimilarityScore > scoredRecords[b].SimilarityScore
	})

	// limit to top k
	if k > len(scoredRecords) {
		k = len(scoredRecords)
	}

	topK := scoredRecords[:k]

	return topK, nil
}
