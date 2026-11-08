// Package flat will implement exact nearest-neighbor search.
package flat

import (
	"errors"
	"fmt"

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
