package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"go.etcd.io/bbolt"
)

func createDb(t *testing.T) (*bbolt.DB, func()) {
	// First, create a temporary directory to be used for the duration of
	// this test.
	tempDirName, err := ioutil.TempDir("", "channeldb")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}
	path := filepath.Join(tempDirName, "testdb.db")

	options := &bbolt.Options{
		NoFreelistSync: true,
		FreelistType:   bbolt.FreelistMapType,
	}

	dbFilePermission := os.FileMode(0600)
	bdb, err := bbolt.Open(path, dbFilePermission, options)
	if err != nil {
		t.Fatalf("error creating bbolt db: %v", err)
	}

	cleanUp := func() {
		bdb.Close()
		os.RemoveAll(tempDirName)
	}

	return bdb, cleanUp
}

func testFunc(t *testing.T) {
	t.Parallel()

	db, cleanUp := createDb(t)
	defer cleanUp()

	targetNb := 100 // You might need to increase this.
	bucketName := []byte("graph-node")

	for i := 0; i < targetNb; i++ {
		err := db.Update(func(tx *bbolt.Tx) error {
			nodes, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				return err
			}

			var key [16]byte
			rand.Read(key[:])
			if err := nodes.Put(key[:], nil); err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestHelloWorld(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("%d", i), testFunc)
	}
}
