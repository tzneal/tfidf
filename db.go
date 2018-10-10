package tfidf

type DB interface {
	DocumentCount() (int, error)
	AddDocument(counts map[string]int) error
	TermOccurrences(text string) (int, error)
}
