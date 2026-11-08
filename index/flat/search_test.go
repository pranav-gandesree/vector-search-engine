package flat

import (
	"math"
	"reflect"
	"testing"

	"vector-search/index"
	"vector-search/vector"
)

func newSearchTestIndex(t *testing.T) *Index {
	t.Helper()
	idx := newTestIndex(t)
	// Insert out of ranking order so search must sort its results.
	for _, record := range []index.Record{
		{Id: "C", Vector: []float64{-1, 0, 0}},
		{Id: "B", Vector: []float64{0, 1, 0}},
		{Id: "D", Vector: []float64{1, 1, 0}},
		{Id: "A", Vector: []float64{1, 0, 0}},
	} {
		if err := idx.Add(record); err != nil {
			t.Fatalf("add %q: %v", record.Id, err)
		}
	}
	return idx
}

func TestSearchRankingAndTopK(t *testing.T) {
	want := []index.SearchResult{
		{ItemId: "A", SimilarityScore: 1},
		{ItemId: "D", SimilarityScore: 1 / math.Sqrt(2)},
		{ItemId: "B", SimilarityScore: 0},
		{ItemId: "C", SimilarityScore: -1},
	}
	for _, tc := range []struct {
		name string
		k    int
		want int
	}{
		{name: "one", k: 1, want: 1},
		{name: "top three", k: 3, want: 3},
		{name: "all", k: 4, want: 4},
		{name: "oversized k", k: 10, want: 4},
	} {
		t.Run(tc.name, func(t *testing.T) {
			idx := newSearchTestIndex(t)
			// A non-unit query verifies that search normalizes its input.
			got, err := idx.Search(vector.Vector{Values: []float64{10, 0, 0}}, tc.k)
			if err != nil {
				t.Fatalf("search: %v", err)
			}
			if len(got) != tc.want {
				t.Fatalf("expected %d results, got %d", tc.want, len(got))
			}
			for i, result := range got {
				if result.ItemId != want[i].ItemId {
					t.Errorf("result %d: expected ID %q, got %q", i, want[i].ItemId, result.ItemId)
				}
				if !(math.Abs(result.SimilarityScore-want[i].SimilarityScore) <= 1e-9) {
					t.Errorf("result %d: expected score %g, got %g", i, want[i].SimilarityScore, result.SimilarityScore)
				}
			}
		})
	}
}

func TestSearchBreaksTiesByID(t *testing.T) {
	idx := newTestIndex(t)
	for _, id := range []string{"z", "m", "a"} {
		if err := idx.Add(index.Record{Id: id, Vector: []float64{1, 0, 0}}); err != nil {
			t.Fatalf("add %q: %v", id, err)
		}
	}
	got, err := idx.Search(vector.Vector{Values: []float64{1, 0, 0}}, 2)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	if got[0].ItemId != "a" || got[1].ItemId != "m" {
		t.Errorf("expected tied IDs [a m], got [%s %s]", got[0].ItemId, got[1].ItemId)
	}
}

func TestSearchEmptyIndex(t *testing.T) {
	idx := newTestIndex(t)
	got, err := idx.Search(vector.Vector{Values: []float64{1, 0, 0}}, 5)
	if err != nil {
		t.Fatalf("search empty index: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected no results, got %v", got)
	}
}

func TestSearchRejectsInvalidRequests(t *testing.T) {
	for _, tc := range []struct {
		name  string
		query []float64
		k     int
	}{
		{name: "zero k", query: []float64{1, 0, 0}, k: 0},
		{name: "negative k", query: []float64{1, 0, 0}, k: -1},
		{name: "wrong dimension", query: []float64{1, 0}, k: 1},
		{name: "empty query", query: nil, k: 1},
		{name: "zero vector", query: []float64{0, 0, 0}, k: 1},
		{name: "NaN", query: []float64{1, math.NaN(), 0}, k: 1},
		{name: "positive infinity", query: []float64{1, math.Inf(1), 0}, k: 1},
		{name: "negative infinity", query: []float64{1, math.Inf(-1), 0}, k: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Invalid queries must fail even when the index is empty.
			for _, populated := range []bool{false, true} {
				idx := newTestIndex(t)
				if populated {
					idx = newSearchTestIndex(t)
				}
				got, err := idx.Search(vector.Vector{Values: tc.query}, tc.k)
				if err == nil {
					t.Errorf("populated=%v: expected an error", populated)
				}
				if got != nil {
					t.Errorf("populated=%v: expected nil results on error, got %v", populated, got)
				}
			}
		})
	}
}

func TestSearchDoesNotMutateInputsOrIndex(t *testing.T) {
	idx := newSearchTestIndex(t)
	query := vector.Vector{Values: []float64{3, 4, 0}}
	originalQuery := append([]float64(nil), query.Values...)
	originalRecords := make([]index.Record, len(idx.Records))
	for i, record := range idx.Records {
		originalRecords[i] = record
		originalRecords[i].Vector = append([]float64(nil), record.Vector...)
	}
	originalIDs := make(map[string]struct{}, len(idx.IDs))
	for id := range idx.IDs {
		originalIDs[id] = struct{}{}
	}
	originalDimension := idx.Dimension
	if _, err := idx.Search(query, 3); err != nil {
		t.Fatalf("search: %v", err)
	}
	if !reflect.DeepEqual(query.Values, originalQuery) {
		t.Errorf("search modified query: expected %v, got %v", originalQuery, query.Values)
	}
	if !reflect.DeepEqual(idx.Records, originalRecords) {
		t.Error("search modified stored records or their order")
	}
	if !reflect.DeepEqual(idx.IDs, originalIDs) {
		t.Error("search modified ID map")
	}
	if idx.Dimension != originalDimension {
		t.Error("search modified index dimension")
	}
}
