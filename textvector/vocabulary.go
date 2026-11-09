package textvector

import "sort"

// it gives all the unique words and what position each word have in a vector
func VocabularyBuilding(documents []string) map[string]int {
	//store unique words
	uniqueWords := make(map[string]bool)

	//tokenize each document
	for _, document := range documents {
		tokens := Tokenize(document)

		for _, token := range tokens {
			uniqueWords[token] = true
		}
	}

	//copy unique words into a slice
	words := make([]string, 0, len(uniqueWords))

	for word := range uniqueWords {
		words = append(words, word)
	}

	//sort words alphabetically
	sort.Strings(words)

	//assign a position to each word
	vocabulary := make(map[string]int)

	for index, word := range words {
		vocabulary[word] = index
	}

	return vocabulary
}
