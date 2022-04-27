package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp		int64
	PrevBlockHash	[]byte
	Hash			[]byte
	Data			[]byte
	//验证工作量证明
	Nonce			int
}

//Serialize 将Block序列化为一个字节数组
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

//DeserializeBlock 将字节数组反序列化为一个Block
func DeserializeBlock(d []byte) *Block {
	var b Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}

	return &b
}

//NewBlock 创建新块时需要运行工作量证明找到有效哈希
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:		time.Now().Unix(),
		PrevBlockHash:	prevBlockHash,
		Hash:			[]byte{},
		Data:			[]byte(data),
		Nonce:			0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}