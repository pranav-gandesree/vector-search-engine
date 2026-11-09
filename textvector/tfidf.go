package textvector

func Transformation(document string, vocabulary map[string]int, IDFWeight map[string]float64) []float64 {
	//calculate term frequecy in each document
	tokens := Tokenize(document)

	result := make(map[string]int)

	for _, token := range tokens {
		result[token]++
	}

	vector := make([]float64, len(vocabulary))

	//calculate TF-IDF for each word
	for word, count := range result {

		// find the word's position
		position, exists := vocabulary[word]

		// skip words not in vocabulary
		if !exists {
			continue
		}

		// retrieve IDF weight
		idf := IDFWeight[word]

		// calculate TF-IDF
		vector[position] = float64(count) * idf
	}

	return vector
}
