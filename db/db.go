package db

import (
	"github.com/boltdb/bolt"
	"github.com/yoonhero/ohpotatocoin/utils"
)

// declare const not to mistake
const (
	dbName       = "blockchain.db"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
)

// bolt.DB pointer
var db *bolt.DB

// create or load database
func DB() *bolt.DB {
	// if db var is nil
	if db == nil {
		// init db
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		db = dbPointer
		utils.HandleErr(err)

		// create bucket if not exist bucket
		err = db.Update(func(t *bolt.Tx) error {
			_, err := t.CreateBucketIfNotExists([]byte(dataBucket))
			utils.HandleErr(err)
			_, err = t.CreateBucketIfNotExists([]byte(blocksBucket))
			return err
		})
		utils.HandleErr(err)
	}

	// return database
	return db
}

// save block data
func SaveBlock(hash string, data []byte) {
	// update database
	err := DB().Update(func(t *bolt.Tx) error {
		// get block bucket
		bucket := t.Bucket([]byte(blocksBucket))

		// put hash (key) and data (value)
		// save data
		err := bucket.Put([]byte(hash), data)
		return err
	})
	utils.HandleErr(err)
}

// close database
func Close() {
	DB().Close()
}

// save block chain
func SaveCheckpoint(data []byte) {
	// update database
	err := DB().Update(func(t *bolt.Tx) error {
		// get blockchain bucket
		bucket := t.Bucket([]byte(dataBucket))

		// put "checkpoint" (key) data (value)
		// save data
		err := bucket.Put([]byte(checkpoint), data)
		return err
	})
	utils.HandleErr(err)
}

// blockchain data in db
func Checkpoint() []byte {
	var data []byte
	// read only func View() to see blockchain
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(dataBucket))
		data = bucket.Get([]byte(checkpoint))
		return nil
	})
	return data
}

// find block in db
func Block(hash string) []byte {
	var data []byte
	// read only func View() to find blocks
	DB().View(func(t *bolt.Tx) error {
		bucket := t.Bucket([]byte(blocksBucket))
		data = bucket.Get([]byte(hash))
		return nil
	})

	return data
}
