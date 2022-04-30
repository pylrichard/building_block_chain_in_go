package main

type Transaction struct {
	Id		[]byte
	In		[]TxInput
	Out 	[]TxOutput
}

type TxInput struct {
	TxId		[]byte
	Out			int
	ScriptSig	string
}

type TxOutput struct {
	Value			int
	ScriptPubKey	string
}