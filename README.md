
tfidf
=========

tfidf is a library for calculating term [frequency-inverse
domain](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) scores across documents.

It differs from some of the other go tfidf libraries in that it uses a store
plugin for persisting corpuses for score computation.  A sample boltdb store is
provided which should serve for most usee cases.

Usage
=====

Constructors are provided so you can use custom ellipsoid parameters, but defaults are
provided for WGS84:

```go
  	model := tfidf.New(tfidf.NewMemoryDB(), tfidf.DefaultOptions())
	gettysburg, _ := ioutil.ReadFile("testdata/gettysburg.txt")
	model.AddDocument(string(gettysburg))
	model.AddDocument("Another document that is in the corpus")
	model.AddDocument("One last document that is inside the corpus")
	query := model.Document("The terms War, Lincoln, and Document are in the Corpus")
	scored, _ := model.ScoredTerms(query)
	for _, term := range scored {
		fmt.Println(term.OriginalTerm, term.Score)
	}
	// Output:
	// lincoln, 0.14362780923945326
	// war, 0.14362780923945326
	// corpus 0.05300875094999672
	// document 0.05300875094999672
```
