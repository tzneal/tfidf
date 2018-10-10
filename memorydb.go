package tfidf

type MemoryDB struct {
	numDocs         uint
	termOccurrences map[string]uint
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		termOccurrences: map[string]uint{},
	}
}

func (m *MemoryDB) AddDocument(counts map[string]uint) error {
	m.numDocs++
	for k := range counts {
		// record that we saw each term
		m.termOccurrences[k] = m.termOccurrences[k] + 1
	}
	return nil
}

func (m *MemoryDB) TermOccurrences(text string) (uint, error) {
	return m.termOccurrences[text], nil
}

func (m *MemoryDB) DocumentCount() (uint, error) {
	return m.numDocs, nil
}
