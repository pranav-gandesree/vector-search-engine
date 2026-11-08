// Package index contains types shared by search index implementations.
package index

type Record struct {
	Id     string
	Text   string
	Vector []float64
}

type SearchResult struct {
	ItemId          string
	SimilarityScore float64
}

// Add a shared index interface once Add and Search have been implemented.
