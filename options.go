package tfidf

import (
	"strings"

	porterstemmer "github.com/blevesearch/go-porterstemmer"
)

// WeightingScheme controls the TFIDF score computation.
type WeightingScheme byte

// WeightingScheme consants.
const (
	// WeightingSchemeOne is tf*idf
	WeightingSchemeOne WeightingScheme = iota
	// WeightingSchemeTwo is 1+log_10(tf)
	WeightingSchemeTwo
	// WeightingSchemeThree is (1+log_10(tf)*idf
	WeightingSchemeThree
)

// Options primarily control how documents are tokenized and processed.
type Options struct {
	// SplitFunc is the function used to split a document into string (e.g. strings.Fields)
	SplitFunc func(text string) []string
	// WeightingScheme controls how the TFIDF score is computed.
	WeightingScheme WeightingScheme
	// SkipPostprocessing is a function, which if set, that can be used to skip
	// postprocessing for a given term.  This is useful for things like hash tags,
	// @usernames, etc. which are recognizable without context and which you want
	// to skip post processing step such as word stemming.
	SkipPostprocessing func(term string) bool
	// Preprocessors are functions applied to the input document that operate
	// across the whole document such as lower casing the entire document.
	Preprocessors []func(text string) string
	// Filters are functions which must return true for each string to ensure that the
	// terms are kept.  This can be used for stopword removal.
	Filters []func(term string) bool
	// Postprocessors are applied to each term after preprocessing, splitting and filtering.
	Postprocessors []func(term string) string
}

// AddPreprocessor adds a pre-processor to be used. This is applied to the text
// before anything else is done, it might be used to remove lowercase the entire
// string, replace binary characters (e.g. Word smart quotes), etc.
func (o *Options) AddPreprocessor(fn func(text string) string) {
	o.Preprocessors = append(o.Preprocessors, fn)
}

// AddFilter adds a filter to be used upon splitting the document text.  This is
// applied to each term after splitting, and each term must pass all of the
// filters or it will be removed.  This is used by DefaultOptions() to remove
// English stopwords.
func (o *Options) AddFilter(fn func(term string) bool) {
	o.Filters = append(o.Filters, fn)
}

// AddPostprocessor adds a post-processor to be used. This is applied to the text
// before anything else is done, it might be used to remove lowercase the entire
// string, replace binary characters (e.g. Word smart quotes), etc.
func (o *Options) AddPostprocessor(fn func(term string) string) {
	o.Postprocessors = append(o.Postprocessors, fn)
}

// MakeRemoveStopwordsFunc returns a function to be used in AddFilter that can
// be used to remove stopwords.
func MakeRemoveStopwordsFunc(list []string) func(term string) bool {
	m := map[string]struct{}{}
	for _, w := range list {
		m[w] = struct{}{}
	}
	return func(term string) bool {
		_, isStopword := m[term]
		return !isStopword
	}
}

// DefaultOptions returns the default options which lowercases the document, removes
// extra spacing, english stopwords as well as performing porter word stemming.
func DefaultOptions() Options {
	o := Options{
		WeightingScheme: WeightingSchemeThree,
		SplitFunc:       strings.Fields,
	}
	o.AddPreprocessor(strings.ToLower)
	o.AddPreprocessor(strings.TrimSpace)
	o.AddFilter(MakeRemoveStopwordsFunc(stopwordsEnglish))
	o.AddPostprocessor(func(term string) string {
		return strings.Trim(term, ",.!?")
	})
	o.AddPostprocessor(func(term string) string {
		runes := porterstemmer.Stem([]rune(term))
		if len(runes) > 0 && runes[len(runes)-1] == '\'' {
			runes = runes[:len(runes)-1]
		}
		return string(runes)

	})
	return o
}
