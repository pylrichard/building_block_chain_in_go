package main

import (
	"encoding/hex"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
)

const dbFile = "block_chain.db"
const blocksBucket = "blocks"
const genesisCoinBaseData = "The times 03/01/2009 chancellor on brink of second bailout for banks"

type BlockChain struct {
	tip []byte
	db	*bolt.DB
}

type BlockChainIterator struct {
	currentHash []byte
	db			*bolt.DB
}

func (bc *BlockChain) MineBlock(txs []*Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(txs, lastHash)

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

func (bci *BlockChainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))
		encodedBlock := bucket.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bci.currentHash = block.PrevBlockHash

	return block
}

func IsDbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func NewBlockChain() *BlockChain {
	if IsDbExists() == false {
		fmt.Println("No existing block chain found, create first one")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

//CreateBlockChain 创建一个新的区块链数据库
//addr 接收挖出创世块的奖励
func CreateBlockChain(addr string) *BlockChain {
	if IsDbExists() {
		fmt.Println("BlockChain already exists")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cbTx := NewCoinBaseTx(addr, genesisCoinBaseData)
		genesis := NewGenesisBlock(cbTx)

		b, err := tx.CreateBucket([]byte(blocksBucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.Hash, genesis.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.Hash)
		if err != nil {
			log.Panic(err)
		}
		tip = genesis.Hash

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := BlockChain{tip, db}

	return &bc
}

//FindUnspentTransactions 找到未花费输出的交易
func (bc *BlockChain) FindUnspentTransactions(addr string) []Transaction {
	var unspentTxs []Transaction
	spendTxs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		b := bci.Next()

		for _, tx := range b.Transactions {
			txId := hex.EncodeToString(tx.Id)

		Output:
			for outId, out := range tx.Out {

			}
		}
	}
}