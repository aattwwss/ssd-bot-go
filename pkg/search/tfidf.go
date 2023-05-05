package search

import (
	"math"
	"strings"
)

type TfIdf struct {
	termFrequencies []map[string]int
	Documents       []string
}

func NewTfIdf(documents []string) *TfIdf {
	tfidf := &TfIdf{
		termFrequencies: make([]map[string]int, len(documents)),
		Documents:       documents,
	}

	for i, doc := range documents {
		terms := strings.Fields(doc)
		tfidf.termFrequencies[i] = make(map[string]int)
		for _, term := range terms {
			tfidf.termFrequencies[i][term]++
		}
	}

	return tfidf
}

func (t *TfIdf) Tf(term string, documentIndex int) float64 {
	freq := t.termFrequencies[documentIndex][term]
	if freq == 0 {
		return 0
	}
	return 1 + math.Log(float64(freq))
}

func (t *TfIdf) Idf(term string) float64 {
	n := float64(len(t.Documents))
	df := 0.0
	for _, doc := range t.Documents {
		if strings.Contains(doc, term) {
			df++
		}
	}
	if df == 0 {
		return 0
	}
	return math.Log(n / df)
}

func (t *TfIdf) TfIdf(term string, documentIndex int) float64 {
	return t.Tf(term, documentIndex) * t.Idf(term)
}

func (t *TfIdf) SearchWithString(searchString string) string {
	terms := strings.Fields(searchString)
	return t.Search(terms)
}

func (t *TfIdf) Search(searchTerms []string) string {
	scores := make([]float64, len(t.Documents))
	for i := range t.Documents {
		for _, term := range searchTerms {
			scores[i] += t.TfIdf(term, i)
		}
	}
	maxScore := 0.0
	maxIndex := 0
	for i, score := range scores {
		if score > maxScore {
			maxScore = score
			maxIndex = i
		}
	}
	return t.Documents[maxIndex]
}
