package block

import (
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/pylrichard/building_block_chain_in_go/simple/transaction"
)

const dbFileNameTemplate = "block_chain_%s.db"
const blocksBucket = "blocks"
const genesisCoinBaseData = "The Times 03/Jan/2009 chancellor on brink of second bailout for banks"

type Chain struct {
	tip []byte
	Db	*bolt.DB
}

func NewChainWithGenesis(addr, nodeId string) *Chain {
	dbFileName := fmt.Sprintf(dbFileNameTemplate, nodeId)
	if IsDbExists(dbFileName) {
		fmt.Println("BlockChain already exists")
		os.Exit(1)
	}

	var tip []byte
	cbTx := transaction.NewCoinBaseTx(addr, genesisCoinBaseData)
	genesis := NewGenesisBlock(cbTx)

	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
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

	bc := Chain{tip, db}

	return &bc
}

func NewChain(nodeId string) *Chain {
	dbFileName := fmt.Sprintf(dbFileNameTemplate, nodeId)
	if IsDbExists(dbFileName) == false {
		fmt.Println("No existing block chain found, Create new one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFileName, 0600, nil)
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

	bc := Chain{tip, db}

	return &bc
}

func IsDbExists(dbFileName string) bool {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		return false
	}

	return true
}

func (bc *Chain) AddBlock(b *Block) {
}

func (bc *Chain) Iterator() *ChainIterator {
	bci := &ChainIterator{bc.tip, bc.Db}

	return bci
}

//GetBestHeight 返回最后一个块的高度
func (bc *Chain) GetBestHeight() int {
	var lastBlock Block

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return lastBlock.Height
}