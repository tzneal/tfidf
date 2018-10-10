package boltdb

import (
	"encoding/binary"
	"fmt"

	"github.com/tzneal/tfidf"
	bbolt "go.etcd.io/bbolt"
)

type BoltStore struct {
	db *bbolt.DB
}

var _ tfidf.Store = (*BoltStore)(nil)

var (
	metaBucket     = []byte("metadata")
	documentBucket = []byte("documents")
	docCountKey    = []byte("documentCount")
)

func NewBoltStore(db *bbolt.DB) (*BoltStore, error) {
	err := db.Update(func(tx *bbolt.Tx) error {
		for _, bkt := range [][]byte{metaBucket, documentBucket} {
			_, err := tx.CreateBucketIfNotExists(bkt)
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &BoltStore{db}, nil
}

func (b *BoltStore) Close() error {
	return b.db.Close()
}

func (b *BoltStore) DocumentCount() (uint, error) {
	var cnt uint
	err := b.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket(metaBucket)
		d := bkt.Get([]byte(docCountKey))
		if d != nil {
			cnt = uint(binary.BigEndian.Uint32(d))
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return cnt, nil
}
func (b *BoltStore) AddDocument(counts map[string]uint) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		meta := tx.Bucket(metaBucket)
		var docCount uint
		d := meta.Get(docCountKey)
		if d != nil {
			docCount = uint(binary.BigEndian.Uint32(d))
		}
		docCount++
		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf[:], uint32(docCount))
		if err := meta.Put(docCountKey, buf[:]); err != nil {
			return err
		}
		docs := tx.Bucket(documentBucket)
		for term := range counts {
			termKey := []byte(term)
			d := docs.Get(termKey)
			var termCnt uint
			if d != nil {
				termCnt = uint(binary.BigEndian.Uint32(d))
			}
			termCnt++
			// can't reuse the buf as it isn't copied until the commit later
			buf = make([]byte, 4)
			binary.BigEndian.PutUint32(buf[:], uint32(termCnt))
			if err := docs.Put(termKey, buf[:]); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
func (b *BoltStore) TermOccurrences(text string) (uint, error) {
	var cnt uint
	err := b.db.View(func(tx *bbolt.Tx) error {
		bkt := tx.Bucket(documentBucket)
		d := bkt.Get([]byte(text))
		if d != nil {
			cnt = uint(binary.BigEndian.Uint32(d))
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func dump(bkt *bbolt.Bucket, name string) {
	bkt.ForEach(func(k, v []byte) error {
		fmt.Println(name, "DUMP", string(k), "=", v)
		return nil
	})
}
