package block

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/pylrichard/building_block_chain_in_go/simple/transaction"
)

type Block struct {
	Timestamp		int64
	Transactions	[]*transaction.Transaction
	PrevBlockHash	[]byte
	Hash			[]byte
	Nonce			int
	Height			int
}

func NewBlock(txs []*transaction.Transaction, prevBlockHash []byte, height int) *Block {
	b := &Block{time.Now().Unix(), txs,
		prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(b)
	nonce, hash := pow.Run()

	b.Hash = hash[:]
	b.Nonce = nonce

	return b
}

func NewGenesisBlock(coinBase *transaction.Transaction) *Block {
	return NewBlock([]*transaction.Transaction{coinBase}, []byte{}, 0)
}

func (b *Block) HashTransaction() []byte {
	var txs [][]byte

	for _, tx := range b.Transactions {
		txs = append(txs, tx.Serialize())
	}

	return txs[1]
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}