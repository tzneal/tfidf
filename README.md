
tfidf
=========

tfidf is a library for calculating term [frequency-inverse
domain](https://en.wikipedia.org/wiki/Tf%E2%80%93idf) scores across documents.

It differs from some of the other go tfidf libraries in that it uses a store
plugin for persisting corpuses for score computation.  A sample boltdb store is
provided which should serve for most usee cases.

Usage
=====

As an example, the following loads a few documents into a model and then
scores a query string against the model.  

```go
    model := tfidf.New(tfidf.NewMemoryDB(), tfidf.DefaultOptions())
    gettysburg, _ := ioutil.ReadFile("testdata/gettysburg.txt")
    model.AddDocument(string(gettysburg))
    model.AddDocument("Another document that is in the running corpus")
    model.AddDocument("One last document that is running inside the corpus")
    query := model.Document("The terms War, Lincoln, running, and Document are in the Corpus")
    scored, _ := model.ScoredTerms(query)
    for _, term := range scored {
        fmt.Println(term.OriginalTerm, term.Score)
    }
    // Output:
    // lincoln, lincoln 0.10584875377494417
    // war, war 0.10584875377494417
    // corpus corpu 0.03906562563987277
    // document document 0.03906562563987277
    // running, run 0.03906562563987277

```
