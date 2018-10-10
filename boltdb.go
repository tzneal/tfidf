package tfidf

import (
	"encoding/binary"
	"fmt"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	db *bolt.DB
}

var _ DB = (*BoltDB)(nil)

var (
	metaBucket     = []byte("metadata")
	documentBucket = []byte("documents")
	docCountKey    = []byte("documentCount")
)

func NewBoltDB(db *bolt.DB) (*BoltDB, error) {
	err := db.Update(func(tx *bolt.Tx) error {
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

// TODO: change to uint32
func (b *BoltDB) DocumentCount() (int, error) {
	var cnt int
	err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(metaBucket)
		d := bkt.Get([]byte(docCountKey))
		if d != nil {
			cnt = int(binary.BigEndian.Uint32(d))
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
func (b *BoltDB) AddDocument(counts map[string]int) error {
	err := b.db.Update(func(tx *bolt.Tx) error {
		meta := tx.Bucket(metaBucket)
		var cnt int
		d := meta.Get(docCountKey)
		if d != nil {
			cnt = int(binary.BigEndian.Uint32(d))
		}
		cnt++
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], uint32(cnt))
		meta.Put(docCountKey, buf[:])

		docs := tx.Bucket(documentBucket)
		for term := range counts {
			termKey := []byte(term)
			d := docs.Get(termKey)
			var cnt int
			if d != nil {
				cnt = int(binary.BigEndian.Uint32(d))
			}
			cnt++
			binary.BigEndian.PutUint32(buf[:], uint32(cnt))
			docs.Put(termKey, buf[:])
		}
		return nil
	})
	return err
}
func (b *BoltDB) TermOccurrences(text string) (int, error) {
	var cnt int
	err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(documentBucket)
		d := bkt.Get([]byte(text))
		if d != nil {
			cnt = int(binary.BigEndian.Uint32(d))
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
