package main

import (
	"fmt"
	blog "vector-search/blogs"
	"vector-search/textvector"
)

func main() {
	// idx, _ := flat.NewIndex(2)

	// record1 := index.Record{
	// 	Id:     "A",
	// 	Text:   "Database indexingg",
	// 	Vector: []float64{1, 0},
	// }

	// record2 := index.Record{
	// 	Id:     "B",
	// 	Text:   "vector search",
	// 	Vector: []float64{0, 1},
	// }

	// record3 := index.Record{
	// 	Id:     "C",
	// 	Text:   "vector search engine",
	// 	Vector: []float64{-1, 0},
	// }

	// record4 := index.Record{
	// 	Id:     "D",
	// 	Text:   "distributed systems",
	// 	Vector: []float64{1, 1},
	// }

	// idx.Add(record1)
	// idx.Add(record2)
	// idx.Add(record3)
	// idx.Add(record4)

	// query := vector.Vector{Values: []float64{0, 1}}

	// results, err := idx.Search(query, 3)

	// if err != nil {
	// 	fmt.Println("Search error:", err)
	// 	return
	// }

	// for i, result := range results {
	// 	fmt.Printf("Search result %d: %+v\n", i+1, result)
	// }

	filePath := "blogs/blogs.json"

	articles, err := blog.LoadArticles(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}

	var documents []string
	for _, article := range articles {
		text := article.Title + " " + article.Body
		documents = append(documents, text)
	}

	fmt.Printf("total documents %d\n", len(documents))

	// build vocabulary from all documents
	vocabulary := textvector.VocabularyBuilding(documents)

	if len(vocabulary) == 0 {
		fmt.Println("Vocabulary is empty")
		return
	}

	fmt.Printf("Vocabulary size: %d\n", len(vocabulary))

	// for word, position := range vocabulary {
	// 	fmt.Printf("Word: %s | Position: %d\n", word, position)
	// }

	documentFrequency := textvector.DocumentFrequency(documents)

	fmt.Println("Document frequencies:")

	// for word, count := range documentFrequency {
	// 	fmt.Printf("%s: %d\n", word, count)
	// }

	idfWeights := textvector.IDF(documentFrequency, len(documents))

	selectedWords := []string{"vector", "database", "ai", "search"}
	for _, word := range selectedWords {
		if weight, exists := idfWeights[word]; exists {
			fmt.Printf("Word: %s | IDF: %.4f\n", word, weight)
		}
	}

	if len(documents) == 0 {
		fmt.Println("No documents available")
		return
	}

	// Transform the first document
	vector := textvector.Transformation(
		documents[0],
		vocabulary,
		idfWeights,
	)

	// Count nonzero components
	nonZeroCount := 0

	for _, value := range vector {
		if value != 0 {
			nonZeroCount++
		}
	}

	// Print summary
	fmt.Printf("Article title: %s\n", articles[0].Title)
	fmt.Printf("Vector length: %d\n", len(vector))
	fmt.Printf("Vocabulary size: %d\n", len(vocabulary))
	fmt.Printf("Nonzero components: %d\n", nonZeroCount)

	// Verify vector length
	if len(vector) != len(vocabulary) {
		fmt.Println("Error: Vector length does not match vocabulary size")
		return
	}

}


func Actual(){
	
}