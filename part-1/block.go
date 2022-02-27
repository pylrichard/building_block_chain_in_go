package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

//Block 由区块头和交易两部分构成，Timestamp/PrevBlockHash/Hash属于区块头
type Block struct {
	//当前时间戳，即区块创建的时间
	Timestamp		int64
	//前一个块的哈希
	PrevBlockHash	[]byte
	//当前块的哈希
	Hash			[]byte
	//区块存储的信息，比特币中就是交易
	Data			[]byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	b := &Block{
		Timestamp:		time.Now().Unix(),
		PrevBlockHash:	prevBlockHash,
		Hash:			[]byte{},
		Data:			[]byte(data),
	}
	b.SetHash()

	return b
}

func (b *Block) SetHash() {
	t := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Data, t}, []byte{})
	h := sha256.Sum256(headers)

	b.Hash = h[:]
}

func NewGenesisBlock() *Block {
	return NewBlock("GenesisBlock", []byte{})
}