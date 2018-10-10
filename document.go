package tfidf

import "sort"

// Document contains term frequency statistics for a single document.
type Document struct {
	total      uint
	counts     map[string]uint
	termmap    map[string]string
	invTermmap map[string]string
}

// TF returns the term frequency for a given term within a document.
func (d Document) TF(term Term) float64 {
	cnt, ok := d.counts[term.ProcessedTerm]
	if !ok || d.total == 0 {
		return 0
	}
	return float64(cnt) / float64(d.total)
}

// Term is a tokenized term from a document. It has both the original form of
// the term as well as the processed version used for lookups.
type Term struct {
	OriginalTerm  string
	ProcessedTerm string
}

// Terms returns the processed terms from the docuoment
func (d Document) Terms() []Term {
	terms := []Term{}
	for v := range d.counts {
		terms = append(terms, Term{
			OriginalTerm:  d.invTermmap[v],
			ProcessedTerm: v,
		})
	}
	// just to get a deterministic order
	sort.Slice(terms, func(a, b int) bool {
		return terms[a].OriginalTerm < terms[b].OriginalTerm
	})
	return terms
}
