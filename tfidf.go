package tfidf

import (
	"math"
)

// Model is used to compute term frequency/inverse document frequency scores.
type Model struct {
	db   Store
	opts Options
}

// New creates a new model with a given database and options.
func New(db Store, opts Options) *Model {
	if opts.SkipPostprocessing == nil {
		opts.SkipPostprocessing = func(string) bool { return false }
	}
	return &Model{db, opts}
}

// Document returns a processed document without adding it to the corpus. This can
// be used for determining relevant terms within a document without affecting future
// scoring.
func (m *Model) Document(text string) Document {
	// preprocess
	for _, pre := range m.opts.Preprocessors {
		text = pre(text)
	}

	terms := m.opts.SplitFunc(text)
	doc := Document{
		counts:     map[string]uint{},
		termmap:    map[string]string{},
		invTermmap: map[string]string{},
	}

	// filter out
	filtered := []string{}
	for _, term := range terms {
		passes := true
		for _, fn := range m.opts.Filters {
			passes = passes && fn(term)
		}
		if !passes {
			continue
		}
		orig := term
		// optionally apply term post-processing
		if !m.opts.SkipPostprocessing(term) {
			for _, fn := range m.opts.Postprocessors {
				term = fn(term)
			}
		}
		// keep original term
		doc.termmap[orig] = term
		doc.invTermmap[term] = orig
		filtered = append(filtered, term)
	}

	terms = filtered
	doc.total = uint(len(terms))
	for _, term := range terms {
		doc.counts[term] = doc.counts[term] + 1
	}
	return doc
}

// AddDocument adds a document to the corpus, returning the processed document.
func (m *Model) AddDocument(text string) (Document, error) {
	doc := m.Document(text)
	return doc, m.db.AddDocument(doc.counts)
}

// IDF returns the inverse document frequency for a given procssed term.
func (m *Model) IDF(term Term) (float64, error) {
	cnt, err := m.db.TermOccurrences(term.ProcessedTerm)
	if err != nil {
		return 0, err
	}
	nDocs, err := m.db.DocumentCount()
	if err != nil {
		return 0, err
	}
	// avoid divide by zero for the no occurence case
	if cnt == 0 {
		return 0, nil
	}
	return math.Log10(float64(nDocs) / float64(cnt)), nil
}

// TFIDF returns the term frequency-inverse document frequency score for
// a given term within the document.
func (m *Model) TFIDF(doc Document, term Term) (float64, error) {
	tf := doc.TF(term)
	idf, err := m.IDF(term)
	if err != nil {
		return 0, err
	}
	switch m.opts.WeightingScheme {
	case WeightingSchemeOne:
		return tf * idf, nil
	case WeightingSchemeTwo:
		return (1 + math.Log10(tf)), nil
	default:
		// do something sane
		fallthrough
	case WeightingSchemeThree:
		return (1 + math.Log10(tf)) * idf, nil
	}
}

// ScoredTerm is a term that has been scored against a corpus.
type ScoredTerm struct {
	Term
	Score float64
}

// ScoredTerms returns the terms from a document that have a non-zero score when
// checked against the corpus.
func (m *Model) ScoredTerms(doc Document) ([]ScoredTerm, error) {
	ret := []ScoredTerm{}
	for _, v := range doc.Terms() {
		score, err := m.TFIDF(doc, v)
		if err != nil {
			return nil, err
		}
		// skip terms that don't contribute
		if score == 0 {
			continue
		}
		st := ScoredTerm{Score: score}
		st.OriginalTerm = v.OriginalTerm
		st.ProcessedTerm = v.ProcessedTerm
		ret = append(ret, st)
	}
	return ret, nil
}
