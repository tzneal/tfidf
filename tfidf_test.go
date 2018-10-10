package tfidf_test

import (
	"io/ioutil"
	"math"
	"os"
	"testing"

	"github.com/tzneal/tfidf"
	"github.com/tzneal/tfidf/boltdb"
	bbolt "go.etcd.io/bbolt"
)

func TestWikipediaExampleMem(t *testing.T) {
	testWikipediaExample(t, tfidf.NewMemoryDB())
}

func TestWikipediaExampleBolt(t *testing.T) {
	td, err := ioutil.TempFile("", "bolttest")
	if err != nil {
		t.Fatalf("error creating temp dir: %s", err)
	}
	td.Close()
	fn := td.Name()
	b, err := bbolt.Open(fn, 0666, nil)
	defer os.Remove(fn)
	if err != nil {
		t.Fatalf("error opening bolt DB: %s", err)
	}
	db, err := boltdb.NewBoltStore(b)
	if err != nil {
		t.Fatalf("error opening bolt DB: %s", err)
	}
	testWikipediaExample(t, db)
}

// pulled from https://en.wikipedia.org/wiki/Tf%E2%80%93idf
func testWikipediaExample(t *testing.T, tdb tfidf.Store) {
	t.Helper()
	defer tdb.Close()
	opts := tfidf.DefaultOptions()
	// Wikipedia example uses this weighting
	opts.WeightingScheme = tfidf.WeightingSchemeOne
	// prevent any pre/post processing that DefaultOptions might perform
	opts.Preprocessors = nil
	opts.Postprocessors = nil
	opts.Filters = nil

	db := tfidf.New(tdb, opts)
	d1, err := db.AddDocument("this is a a sample")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	d2, err := db.AddDocument("this is another another example example example")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if got := d1.TF(tfidf.Term{ProcessedTerm: "this"}); got != 0.2 {
		t.Errorf("expected 0.2, got %f", got)
	}
	if got := d2.TF(tfidf.Term{ProcessedTerm: "this"}); math.Abs(got-0.142857) > 1e-5 {
		t.Errorf("expected 0.142857, got %f", got)
	}

	idf, err := db.IDF(tfidf.Term{ProcessedTerm: "this"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if idf != 0 {
		t.Errorf("expected 0, got %f", idf)
	}

	tfsc, err := db.TFIDF(d1, tfidf.Term{ProcessedTerm: "example"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if tfsc != 0 {
		t.Errorf("expected 0, got %f", tfsc)
	}

	tfsc, err = db.TFIDF(d2, tfidf.Term{ProcessedTerm: "example"})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if math.Abs(tfsc-0.129013) > 1e-5 {
		t.Errorf("expected 0.129013, got %f", tfsc)
	}
}

func TestTweets(t *testing.T) {
	opts := tfidf.DefaultOptions()
	opts.SkipPostprocessing = func(term string) bool {
		return term[0] == '#' || term[0] == '@'
	}
	db := tfidf.New(tfidf.NewMemoryDB(), opts)
	doc := db.Document("This is some #hashtaggy #text from @johndoe")
	// the sentence is just stop words, hashtags and usernames.  The
	// skip post processing function tells us to not post-process the hashtags
	// or username, so their term post-processed form should match the original
	// term
	for _, term := range doc.Terms() {
		if term.OriginalTerm != term.ProcessedTerm {
			t.Errorf("expected no change, got %v", term)
		}
	}
}

func TestGettysburg(t *testing.T) {
	db := tfidf.New(tfidf.NewMemoryDB(), tfidf.DefaultOptions())
	gettysburg, err := ioutil.ReadFile("testdata/gettysburg.txt")
	if err != nil {
		t.Fatalf("unable to open test file: %s", err)
	}
	for i := 0; i < 1000; i++ {
		switch i % 100 {
		case 0:
			db.AddDocument(string(gettysburg) + " quux")
		case 1:
			db.AddDocument(string(gettysburg) + " frobnitz")
		default:
			db.AddDocument(string(gettysburg))
		}

	}
	doc := db.Document("We have come to FROBnitz's dedicate a portion of that field quuxes")
	scored, err := db.ScoredTerms(doc)
	if err != nil {
		t.Fatalf("error scoring: %s", err)
	}
	// All the other terms in our search document appeared in every document in
	// the corpus except for 'quux' and 'frobnitz' which only occurred in a
	// subset.  This means they are the only terms that are of interest when
	// searching (since every other term matches all documents)
	if len(scored) != 2 {
		t.Errorf("expected two scored terms, found %d (%v)", len(scored), scored)
	}
}

func BenchmarkGettysburg10(b *testing.B) { benchmarkGettysburg(b, .1) }
func BenchmarkGettysburg20(b *testing.B) { benchmarkGettysburg(b, .2) }
func BenchmarkGettysburg30(b *testing.B) { benchmarkGettysburg(b, .3) }
func BenchmarkGettysburg40(b *testing.B) { benchmarkGettysburg(b, .4) }
func BenchmarkGettysburg50(b *testing.B) { benchmarkGettysburg(b, .5) }

func benchmarkGettysburg(b *testing.B, pct float64) {
	b.Helper()
	opts := tfidf.DefaultOptions()
	opts.Postprocessors = nil
	opts.Preprocessors = nil
	opts.Filters = nil
	db := tfidf.New(tfidf.NewMemoryDB(), opts)
	gettysburg, err := ioutil.ReadFile("testdata/gettysburg.txt")
	if err != nil {
		b.Fatalf("unable to open test file: %s", err)
	}
	gbstring := string(gettysburg)
	gbstring = gbstring[:int(float64(len(gbstring))*pct)]
	for i := 0; i < b.N; i++ {
		db.AddDocument(gbstring)
	}
}
