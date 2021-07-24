package db

import (
	"fmt"
	"os"

	"github.com/yoonhero/ohpotatocoin/utils"
	bolt "go.etcd.io/bbolt"
)

// declare const not to mistake
const (
	dbName       = "blockchain"
	dataBucket   = "data"
	blocksBucket = "blocks"
	checkpoint   = "checkpoint"
	blockChainDb = "blockchain.db"
)

// bolt.DB pointer
var db *bolt.DB

func getDbName() string {
	// port := os.Args[2][6:]
	port := os.Getenv("PORT")
	return fmt.Sprintf("%s_%s.db", dbName, port)
}

// create or load database
func DB() *bolt.DB {
	// if db var is nil
	if db == nil {
		// init db
		dbPointer, err := bolt.Open(blockChainDb, 0600, nil)
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

func EmptyBlocks() {
	DB().Update(func(t *bolt.Tx) error {
		utils.HandleErr(t.DeleteBucket([]byte(blocksBucket)))
		_, err := t.CreateBucket([]byte(blocksBucket))
		utils.HandleErr(err)
		return nil
	})
}
