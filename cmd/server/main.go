package main

import (
	"fmt"
	"os"
	"strings"

	blog "vector-search/blogs"
	"vector-search/index"
	"vector-search/index/flat"
	"vector-search/textvector"
	"vector-search/vector"
)

func main() {
	filePath := "blogs/blogs.json"

	articles, err := blog.LoadArticles(filePath)

	if err != nil {
		fmt.Println(err)
		return
	}
	if len(articles) == 0 {
		fmt.Println("No documents available")
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
	articleVector := textvector.Transformation(
		documents[0],
		vocabulary,
		idfWeights,
	)

	// Count nonzero components
	nonZeroCount := 0

	for _, value := range articleVector {
		if value != 0 {
			nonZeroCount++
		}
	}

	// Print summary
	fmt.Printf("Article title: %s\n", articles[0].Title)
	fmt.Printf("Vector length: %d\n", len(articleVector))
	fmt.Printf("Vocabulary size: %d\n", len(vocabulary))
	fmt.Printf("Nonzero components: %d\n", nonZeroCount)

	// Verify vector length
	if len(articleVector) != len(vocabulary) {
		fmt.Println("Error: Vector length does not match vocabulary size")
		return
	}

	// index every article using the same vocabulary and IDF weights
	idx, err := flat.NewIndex(len(vocabulary))
	if err != nil {
		fmt.Println("Unable to create index:", err)
		return
	}
	articlesByID := make(map[string]blog.Article, len(articles))
	for position, article := range articles {
		values := textvector.Transformation(documents[position], vocabulary, idfWeights)
		if err := idx.Add(index.Record{
			Id: article.ID, Text: documents[position], Vector: values,
		}); err != nil {
			fmt.Printf("Unable to index article %q: %v\n", article.ID, err)
			return
		}
		articlesByID[article.ID] = article
	}
	fmt.Printf("Indexed articles: %d\n", len(articlesByID))

	queryText := "virtual memory page tables"
	if len(os.Args) > 1 {
		queryText = strings.Join(os.Args[1:], " ")
	}
	fmt.Printf("\nQuery: %s\n", queryText)
	queryValues := textvector.Transformation(queryText, vocabulary, idfWeights)
	hasKnownWords := false
	for _, value := range queryValues {
		if value != 0 {
			hasKnownWords = true
			break
		}
	}
	if !hasKnownWords {
		fmt.Println("No matching vocabulary")
		return
	}

	results, err := idx.Search(vector.Vector{Values: queryValues}, 3)
	if err != nil {
		fmt.Println("Search error:", err)
		return
	}
	shown := 0
	for _, result := range results {
		if result.SimilarityScore <= 0 {
			continue
		}
		article := articlesByID[result.ItemId]
		shown++
		fmt.Printf("%d. %s | Score: %.4f\n   %s\n", shown, article.Title, result.SimilarityScore, article.URL)
	}
	if shown == 0 {
		fmt.Println("No matching articles")
	}
}
