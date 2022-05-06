package cli

import (
	"log"

	"github.com/pylrichard/building_block_chain_in_go/simple/block"
	"github.com/pylrichard/building_block_chain_in_go/simple/wallet"
)

func (cli *CLI) createBlockChain(addr, nodeId string) {
	if !wallet.ValidateAddr(addr) {
		log.Panic("Error: addr is not valid")
	}
	bc := block.CreateChain(addr, nodeId)
	defer bc.Db.Close()
}