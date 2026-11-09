package blog

import (
	"encoding/json"
	"os"
)

type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Body  string `json:"body"`
}

func LoadArticles(file string) ([]Article, error) {
	data, err := os.ReadFile(file)

	if err != nil {
		return nil, err
	}

	var articles []Article

	err = json.Unmarshal(data, &articles)

	if err != nil {
		return nil, err
	}

	return articles, nil

}
