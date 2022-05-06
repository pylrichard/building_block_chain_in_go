package transaction

import (
	"bytes"

	"github.com/pylrichard/building_block_chain_in_go/simple/wallet"
)

type TxInput struct {
	TxId		[]byte
	Out			int
	Signature	[]byte
	PubKey		[]byte
}

func (in *TxInput) IsKeyUsed(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}