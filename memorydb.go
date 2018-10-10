package tfidf

type MemoryDB struct {
	numDocs         int
	termOccurrences map[string]int
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		termOccurrences: map[string]int{},
	}
}

func (m *MemoryDB) AddDocument(counts map[string]int) error {
	m.numDocs++
	for k := range counts {
		// record that we saw each term
		m.termOccurrences[k] = m.termOccurrences[k] + 1
	}
	return nil
}

func (m *MemoryDB) TermOccurrences(text string) (int, error) {
	return m.termOccurrences[text], nil
}

func (m *MemoryDB) DocumentCount() (int, error) {
	return m.numDocs, nil
}
