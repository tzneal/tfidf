package tfidf

type DB interface {
	DocumentCount() (uint, error)
	AddDocument(counts map[string]uint) error
	TermOccurrences(text string) (uint, error)
}
