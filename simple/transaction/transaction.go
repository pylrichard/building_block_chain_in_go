package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
)

const subsidy = 10

type Transaction struct {
	Id	[]byte
	In	[]TxInput
	Out []TxOutput
}

func NewCoinBaseTx(to, data string) *Transaction {
	if data == "" {
		randData := make([]byte, 20)
		_, err := rand.Read(randData)
		if err != nil {
			log.Panic(err)
		}

		data = fmt.Sprintf("%x", randData)
	}

	txIn := TxInput{[]byte{}, -1, nil, []byte(data)}
	txOut := NewTxOutput(subsidy, to)
	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{*txOut}}
	tx.Id = tx.Hash()

	return &tx
}

func DeserializeTransaction(data []byte) Transaction {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Panic(err)
	}

	return tx
}

func (tx Transaction) IsCoinBase() bool {
	return len(tx.In) == 1 && len(tx.In[0].TxId) == 0 && tx.In[0].Out == -1
}

func (tx Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	return encoded.Bytes()
}

func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.Id = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

//TrimmedCopy 创建Transaction的副本，Signing需要使用
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	for _, input := range tx.In {
		inputs = append(inputs, TxInput{input.TxId, input.Out, nil, nil})
	}

	for _, output := range tx.Out {
		outputs = append(outputs, TxOutput{output.Value, output.PubKeyHash})
	}

	txCopy := Transaction{tx.Id, inputs, outputs}

	return txCopy
}

func (tx *Transaction) Verify(prevTxs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}

	for _, input := range tx.In {
		if prevTxs[hex.EncodeToString(input.TxId)].Id == nil {
			log.Panic("Error: Previous transaction is not correct")
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for id, input := range tx.In {
		prevTx := prevTxs[hex.EncodeToString(input.TxId)]
		txCopy.In[id].Signature = nil
		txCopy.In[id].PubKey = prevTx.Out[input.Out].PubKeyHash

		r := big.Int{}
		s := big.Int{}
		sigLen := len(input.Signature)
		r.SetBytes(input.Signature[:(sigLen / 2)])
		s.SetBytes(input.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(input.PubKey)
		x.SetBytes(input.PubKey[:(keyLen / 2)])
		y.SetBytes(input.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}

		txCopy.In[id].PubKey = nil
	}

	return true
}