package block

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

type ChainIterator struct {
	currentHash []byte
	db			*bolt.DB
}

func (ci *ChainIterator) Next() *Block {
	var block *Block

	err := ci.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(ci.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	ci.currentHash = block.PrevBlockHash

	return block
}