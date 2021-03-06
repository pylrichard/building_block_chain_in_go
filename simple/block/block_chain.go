package block

import (
	"bytes"
	"encoding/hex"
	"errors"
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

func (bc *Chain) AddBlock(b *Block) {
}

func (bc *Chain) FindTransaction(Id []byte) (transaction.Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.Id, Id) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return transaction.Transaction{}, errors.New("transaction is not found")
}

func (bc *Chain) Iterator() *ChainIterator {
	bci := &ChainIterator{bc.tip, bc.Db}

	return bci
}

//GetBestHeight ??????????????????????????????
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

//GetBlock ????????????????????????
func (bc *Chain) GetBlock(blockHash []byte) (Block, error) {
	var block Block

	err := bc.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blocksBucket))

		blockData := bucket.Get(blockHash)
		if blockData == nil {
			return errors.New("block is not found")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return block, err
	}

	return block, nil
}

//GetBlockHashes ???????????????????????????????????????
func (bc *Chain) GetBlockHashes() [][]byte {
	var blocks [][]byte
	bci := bc.Iterator()

	for {
		b := bci.Next()
		blocks = append(blocks, b.Hash)

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

	return blocks
}

func (bc *Chain) MineBlock(transactions []*transaction.Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if bc.VerifyTransaction(tx) != true {
			log.Panic("Error: Invalid transaction")
		}
	}

	err := bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte("l"))

		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)

		lastHeight = block.Height

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(transactions, lastHash, lastHeight + 1)

	err = bc.Db.Update(func(tx *bolt.Tx) error {
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
	if err != nil {
		log.Panic(err)
	}

	return newBlock
}

//VerifyTransaction ??????Transaction???Input Signatures
func (bc *Chain) VerifyTransaction(tx *transaction.Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	prevTxs := make(map[string]transaction.Transaction)
	for _, input := range tx.In {
		prevTx, err := bc.FindTransaction(input.TxId)
		if err != nil {
			log.Panic(err)
		}
		prevTxs[hex.EncodeToString(prevTx.Id)] = prevTx
	}

	return tx.Verify(prevTxs)
}

func IsDbExists(dbFileName string) bool {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		return false
	}

	return true
}