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
	spentTxs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		b := bci.Next()

		for _, tx := range b.Transactions {
			txId := hex.EncodeToString(tx.Id)

		Outputs:
			for outIdx, out := range tx.Out {
				//如果交易输出被花费
				if spentTxs[txId] != nil {
					for _, spentOut := range spentTxs[txId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				//如果该交易输出可以被解锁，即可被花费
				if out.CanBeUnlockedWith(addr) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			if tx.IsCoinBase() == false {
				for _, in := range tx.In {
					if in.CanUnlockOutputWith(addr) {
						txId := hex.EncodeToString(in.TxId)
						spentTxs[txId] = append(spentTxs[txId], in.Out)
					}
				}
			}
		}

		if len(b.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTxs
}

func (bc *BlockChain) FindUTX(addr string) []TxOutput {
	var txOutputs []TxOutput
	unspentTxs := bc.FindUnspentTransactions(addr)

	for _, tx := range unspentTxs {
		for _, out := range tx.Out {
			if out.CanBeUnlockedWith(addr) {
				txOutputs = append(txOutputs, out)
			}
		}
	}

	return txOutputs
}

//FindSpendableOutputs 从addr中找到至少amount的UTXO
func (bc *BlockChain) FindSpendableOutputs(addr string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTxs := bc.FindUnspentTransactions(addr)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txId := hex.EncodeToString(tx.Id)

		for outIdx, out := range tx.Out {
			if out.CanBeUnlockedWith(addr) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txId] = append(unspentOutputs[txId], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}