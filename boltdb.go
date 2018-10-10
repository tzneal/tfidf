package tfidf

import (
	"encoding/binary"
	"fmt"

	bbolt "go.etcd.io/bbolt"
)

type BoltDB struct {
	db *bbolt.DB
}

var _ DB = (*BoltDB)(nil)

var (
	metaBucket     = []byte("metadata")
	documentBucket = []byte("documents")
	docCountKey    = []byte("documentCount")
)

func NewBoltDB(db *bbolt.DB) (*BoltDB, error) {
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
	return &BoltDB{db}, nil
}
func (b *BoltDB) Close() error {
	return b.db.Close()
}

// TODO: change to uuint32
func (b *BoltDB) DocumentCount() (uint, error) {
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
func (b *BoltDB) AddDocument(counts map[string]uint) error {
	err := b.db.Update(func(tx *bbolt.Tx) error {
		meta := tx.Bucket(metaBucket)
		var cnt uint
		d := meta.Get(docCountKey)
		if d != nil {
			cnt = uint(binary.BigEndian.Uint32(d))
		}
		cnt++
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], uint32(cnt))
		meta.Put(docCountKey, buf[:])

		docs := tx.Bucket(documentBucket)
		for term := range counts {
			termKey := []byte(term)
			d := docs.Get(termKey)
			var cnt uint
			if d != nil {
				cnt = uint(binary.BigEndian.Uint32(d))
			}
			cnt++
			binary.BigEndian.PutUint32(buf[:], uint32(cnt))
			docs.Put(termKey, buf[:])
		}
		return nil
	})
	return err
}
func (b *BoltDB) TermOccurrences(text string) (uint, error) {
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
