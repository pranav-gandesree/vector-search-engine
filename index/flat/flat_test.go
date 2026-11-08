package flat

import (
	"math"
	"reflect"
	"testing"

	"vector-search/index"
)

func newTestIndex(t *testing.T) *Index {
	t.Helper()
	idx, err := NewIndex(3)
	if err != nil {
		t.Fatalf("create index: %v", err)
	}
	return idx
}

func assertVectorClose(t *testing.T, got, want []float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected %d components, got %d", len(want), len(got))
	}
	for i := range want {
		if !(math.Abs(got[i]-want[i]) <= 1e-9) {
			t.Errorf("component %d: expected %g, got %g", i, want[i], got[i])
		}
	}
}

func assertEmptyIndex(t *testing.T, idx *Index) {
	t.Helper()
	if len(idx.Records) != 0 {
		t.Errorf("expected no records, got %d", len(idx.Records))
	}
	if len(idx.IDs) != 0 {
		t.Errorf("expected no registered IDs, got %d", len(idx.IDs))
	}
}

func TestAddSuccess(t *testing.T) {
	idx := newTestIndex(t)
	record := index.Record{
		Id:     "doc1",
		Text:   "this is about database indexing",
		Vector: []float64{3, 4, 0},
	}
	if err := idx.Add(record); err != nil {
		t.Fatalf("add record: %v", err)
	}
	if len(idx.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(idx.Records))
	}
	stored := idx.Records[0]
	if stored.Id != record.Id {
		t.Errorf("expected ID %q, got %q", record.Id, stored.Id)
	}
	if stored.Text != record.Text {
		t.Errorf("expected text %q, got %q", record.Text, stored.Text)
	}
	assertVectorClose(t, stored.Vector, []float64{0.6, 0.8, 0})
	if _, exists := idx.IDs[record.Id]; !exists {
		t.Errorf("expected %q in IDs map", record.Id)
	}
	if len(idx.IDs) != 1 {
		t.Errorf("expected 1 registered ID, got %d", len(idx.IDs))
	}
}

func TestAddCopiesInput(t *testing.T) {
	idx := newTestIndex(t)
	input := []float64{3, 4, 0}
	original := append([]float64(nil), input...)
	if err := idx.Add(index.Record{Id: "doc1", Vector: input}); err != nil {
		t.Fatalf("add record: %v", err)
	}
	if !reflect.DeepEqual(input, original) {
		t.Errorf("insertion modified input: expected %v, got %v", original, input)
	}
	if len(idx.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(idx.Records))
	}
	input[0] = 999
	assertVectorClose(t, idx.Records[0].Vector, []float64{0.6, 0.8, 0})
}

func TestAddRejectsDuplicateID(t *testing.T) {
	idx := newTestIndex(t)
	if err := idx.Add(index.Record{
		Id: "doc1", Text: "original article", Vector: []float64{3, 4, 0},
	}); err != nil {
		t.Fatalf("add original record: %v", err)
	}
	if len(idx.Records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(idx.Records))
	}
	original := idx.Records[0]
	original.Vector = append([]float64(nil), original.Vector...)
	if err := idx.Add(index.Record{
		Id: "doc1", Text: "replacement article", Vector: []float64{0, 0, 1},
	}); err == nil {
		t.Error("expected an error for a duplicate ID")
	}
	if len(idx.Records) != 1 {
		t.Fatalf("expected 1 record after rejection, got %d", len(idx.Records))
	}
	if !reflect.DeepEqual(idx.Records[0], original) {
		t.Errorf("duplicate insertion changed original: expected %+v, got %+v", original, idx.Records[0])
	}
	if len(idx.IDs) != 1 {
		t.Errorf("expected 1 registered ID, got %d", len(idx.IDs))
	}
	if _, exists := idx.IDs["doc1"]; !exists {
		t.Error("original ID missing after rejection")
	}
}

func TestAddRejectsWrongDimension(t *testing.T) {
	idx := newTestIndex(t)
	if err := idx.Add(index.Record{Id: "doc1", Vector: []float64{3, 4}}); err == nil {
		t.Error("expected an error for incorrect dimensions")
	}
	assertEmptyIndex(t, idx)
}

func TestAddRejectsZeroVector(t *testing.T) {
	idx := newTestIndex(t)
	if err := idx.Add(index.Record{Id: "doc1", Vector: []float64{0, 0, 0}}); err == nil {
		t.Error("expected an error for a zero vector")
	}
	assertEmptyIndex(t, idx)
}

func TestAddRejectsEmptyID(t *testing.T) {
	idx := newTestIndex(t)
	if err := idx.Add(index.Record{Vector: []float64{3, 4, 0}}); err == nil {
		t.Error("expected an error for an empty ID")
	}
	assertEmptyIndex(t, idx)
}
