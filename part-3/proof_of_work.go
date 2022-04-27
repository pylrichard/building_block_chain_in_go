package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

//targetBits 难度值，哈希前24位必须是0
const targetBits = 24

const maxNonce = math.MaxInt64

//ProofOfWork 每个块的工作量必须证明，都有一个指向Block的指针
//target是目标，最终要找到的哈希必须小于目标
type ProofOfWork struct {
	block	*Block
	target	*big.Int
}

//NewProofOfWork 注意target计算方法
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	//1左移256 - targetBits位
	target.Lsh(target, uint(256 - targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

//工作量证明用到的数据有: PrevBlockHash, Data, Timestamp, targetBits, nonce
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

//Run 工作量证明的核心是寻找有效哈希
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash	[32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("\r%x", hash)
			break
		} else {
			nonce++
		}
	}
	fmt.Print("\n\n")

	return nonce, hash[:]
}

//Validate 验证工作量，只要哈希小于目标就是有效工作量
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}