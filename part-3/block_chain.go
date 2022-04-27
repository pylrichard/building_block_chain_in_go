package main

import (
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
)

const dbFile = "block_chain.db"
const blocksBucket = "blocks"

type BlockChain struct {
	//tip存储最后一个块的哈希。在链末端可能出现短暂分叉，tip代表选择了哪条链
	tip []byte
	//db存储数据库连接
	db	*bolt.DB
}

func NewBlockChain() *BlockChain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		//如果数据库中不存在区块链就创建一个，否则直接读取最后一个块的哈希
		if b == nil {
			fmt.Println("No existing block_chain found. Creating a new one...")
			g := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(g.Hash, g.Serialize())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), g.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = g.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

//AddBlock 加入区块时需要将区块持久化到数据库
func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)
	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash

		return nil
	})
}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci := &BlockChainIterator{bc.tip, bc.db}

	return bci
}

type BlockChainIterator struct {
	currentHash []byte
	db			*bolt.DB
}

//Next 返回链中下一个块
func (bci *BlockChainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bci.currentHash = block.PrevBlockHash

	return block
}