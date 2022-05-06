package transaction

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/pylrichard/building_block_chain_in_go/simple/codec"
)

type TxOutput struct {
	Value		int
	PubKeyHash	[]byte
}

func (out *TxOutput) Lock(addr []byte) {
	pubKeyHash := codec.Base58Decode(addr)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - 4]
	out.PubKeyHash = pubKeyHash
}

func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTxOutput(value int, addr string) *TxOutput {
	output := &TxOutput{value, nil}
	output.Lock([]byte(addr))

	return output
}

type TxOutputs struct {
	Outputs []TxOutput
}

func (outs TxOutputs) Serialize() []byte {
	var buff bytes.Buffer

	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TxOutputs {
	var outputs TxOutputs

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}

	return outputs
}