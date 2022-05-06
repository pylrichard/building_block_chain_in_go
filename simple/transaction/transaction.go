package transaction

type Transaction struct {
	Id	[]byte
	In	[]TxInput
	Out []TxOutput
}

func NewCoinBaseTx(to, data string) *Transaction {

}