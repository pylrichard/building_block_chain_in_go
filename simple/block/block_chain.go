package block

import (
	"fmt"
	"os"

	bolt "go.etcd.io/bbolt"

	"github.com/pylrichard/building_block_chain_in_go/simple/transaction"
)

const dbFileNamePattern = "block_chain_%s.db"
const blocksBucket = "blocks"
const genesisCoinBaseData = "The Times 03/Jan/2009 chancellor on brink of second bailout for banks"

type Chain struct {
	tip []byte
	Db	*bolt.DB
}

func CreateChain(addr, nodeId string) *Chain {
	dbFileName := fmt.Sprintf(dbFileNamePattern, nodeId)
	if IsDbExists(dbFileName) {
		fmt.Println("BlockChain already exists")
		os.Exit(1)
	}

	var tip []byte
	cbTx := transaction.NewCoinBaseTx(addr, genesisCoinBaseData)
}

func IsDbExists(dbFileName string) bool {
	if _, err := os.Stat(dbFileName); os.IsNotExist(err) {
		return false
	}

	return true
}