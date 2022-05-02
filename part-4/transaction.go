package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

type Transaction struct {
	Id  []byte
	In  []TxInput
	Out []TxOutput
}

type TxInput struct {
	TxId      []byte
	Out       int
	ScriptSig string
}

type TxOutput struct {
	Value        int
	ScriptPubKey string
}

func (tx Transaction) IsCoinBase() bool {
	return len(tx.In) == 1 && len(tx.In[0].TxId) == 0 && tx.In[0].Out == -1
}

func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte

	e := gob.NewEncoder(&encoded)
	err := e.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.Id = hash[:]
}

func (in *TxInput) CanUnlockOutputWith(data string) bool {
	return in.ScriptSig == data
}

func (out *TxOutput) CanBeUnlockedWith(data string) bool {
	return out.ScriptPubKey == data
}

func NewCoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{subsidy, to}
	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetId()

	return &tx
}

func NewUTXOTransaction(from, to string, amount int, bc *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("No enough funds")
	}

	for idx, outs := range validOutputs {
		txId, err := hex.DecodeString(idx)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TxInput{txId, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetId()

	return &tx
}
