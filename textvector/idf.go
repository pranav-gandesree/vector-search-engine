package textvector

import "math"

// calculate how many documents contain each word, not how many times each word appear
func DocumentFrequency(documents []string) map[string]int {

	result := make(map[string]int)

	for _, document := range documents {
		tokens := Tokenize(document)

		//fresh set for each document
		seen := make(map[string]bool)

		for _, token := range tokens {
			if !seen[token] {
				result[token]++
				seen[token] = true
			}
		}
	}

	return result
}

// TF (Term Frequency): How many times a word appears in a document
// IDF (Inverse Document Frequency): How rare or common a word is across all documents
func IDF(documentFrequencyMap map[string]int, count int) map[string]float64 {
	result := make(map[string]float64)

	for word, df := range documentFrequencyMap {
		idf := math.Log(
			float64(count+1)/float64(df+1),
		) + 1

		result[word] = idf
	}

	return result
}
